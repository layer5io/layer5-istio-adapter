// nolint:lll
// Generates the layer5 adapter's resource yaml. It contains the adapter's configuration, name,
// supported template names (metric in this case), and whether it is session or no-session based.
//go:generate $REPO_ROOT/bin/mixer_codegen.sh -a mixer/adapter/layer5/config/config.proto -x "-s=false -n layer5 -t metric"
package layer5

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"bytes"
	"os"

	"istio.io/api/mixer/adapter/model/v1beta1"
	policy "istio.io/api/policy/v1beta1"
	"istio.io/istio/mixer/adapter/layer5/config"
	"istio.io/istio/mixer/pkg/adapter"
	"istio.io/istio/mixer/pkg/runtime/handler"
	"istio.io/istio/mixer/template/metric"
	"istio.io/pkg/pool"
)

type (
	// Server is basic server interface
	Server interface {
		Addr() string
		Close() error
		Run(shutdown chan error)
	}

	// Layer5Adapter supports metric template.
	Layer5Adapter struct {
		listener net.Listener
		server   *grpc.Server
		env      adapter.Env
	}
)

var _ metric.HandleMetricServiceServer = &Layer5Adapter{}

// HandleMetric records metric entries
func (s *Layer5Adapter) HandleMetric(ctx context.Context, r *metric.HandleMetricRequest) (*v1beta1.ReportResult, error) {

	s.env.Logger().Infof("received request %v\n", *r)
	var b bytes.Buffer
	cfg := &config.Params{}

	if r.AdapterConfig != nil {
		if err := cfg.Unmarshal(r.AdapterConfig.Value); err != nil {
			s.env.Logger().Errorf("error unmarshalling adapter config: %v", err)
			return nil, err
		}
	}

	b.WriteString(fmt.Sprintf("HandleMetric invoked with:\n  Adapter config: %s\n  Instances: %s\n",
		cfg.String(), instances(r.Instances)))

	if cfg.FilePath == "" {
		s.env.Logger().Infof("%s", b.String())
	} else {
		_, err := os.OpenFile("out.txt", os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			s.env.Logger().Errorf("error creating file: %v", err)
		}
		f, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			s.env.Logger().Errorf("error opening file for append: %v", err)
		}

		defer f.Close()

		s.env.Logger().Infof("writing instances to file %s", f.Name())
		if _, err = f.Write(b.Bytes()); err != nil {
			s.env.Logger().Errorf("error writing to file: %v", err)
		}
	}

	s.env.Logger().Infof("success!!")
	return &v1beta1.ReportResult{}, nil
}

func decodeDimensions(in map[string]*policy.Value) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = decodeValue(v.GetValue())
	}
	return out
}

func decodeValue(in interface{}) interface{} {
	switch t := in.(type) {
	case *policy.Value_StringValue:
		return t.StringValue
	case *policy.Value_Int64Value:
		return t.Int64Value
	case *policy.Value_DoubleValue:
		return t.DoubleValue
	default:
		return fmt.Sprintf("%v", in)
	}
}

func instances(in []*metric.InstanceMsg) string {
	var b bytes.Buffer
	for _, inst := range in {
		b.WriteString(fmt.Sprintf("'%s':\n"+
			"  {\n"+
			"		Value = %v\n"+
			"		Dimensions = %v\n"+
			"  }", inst.Name, decodeValue(inst.Value.GetValue()), decodeDimensions(inst.Dimensions)))
	}
	return b.String()
}

// Addr returns the listening address of the server
func (s *Layer5Adapter) Addr() string {
	return s.listener.Addr().String()
}

// Run starts the server run
func (s *Layer5Adapter) Run(shutdown chan error) {
	shutdown <- s.server.Serve(s.listener)
}

// Close gracefully shuts down the server; used for testing
func (s *Layer5Adapter) Close() error {
	if s.server != nil {
		s.server.GracefulStop()
	}

	if s.listener != nil {
		_ = s.listener.Close()
	}

	return nil
}

func getServerTLSOption(credential, privateKey, caCertificate string) (grpc.ServerOption, error) {
	certificate, err := tls.LoadX509KeyPair(
		credential,
		privateKey,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load key cert pair")
	}
	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile(caCertificate)
	if err != nil {
		return nil, fmt.Errorf("failed to read client ca cert: %s", err)
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		return nil, fmt.Errorf("failed to append client certs")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	}
	tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert

	return grpc.Creds(credentials.NewTLS(tlsConfig)), nil
}

// NewLayer5Adapter creates a new IBP adapter that listens at provided port.
func NewLayer5Adapter(addr string) (Server, error) {
	gp := pool.NewGoroutinePool(5, false)
	env := handler.NewEnv(0, "layer5", gp)
	if addr == "" {
		addr = "0"
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", addr))
	if err != nil {
		return nil, env.Logger().Errorf("unable to listen on socket: %v", err)
	}

	s := &Layer5Adapter{
		listener: listener,
		env:      env,
	}
	env.Logger().Infof("listening on \"%v\"", s.Addr())

	credential := os.Getenv("GRPC_ADAPTER_CREDENTIAL")
	privateKey := os.Getenv("GRPC_ADAPTER_PRIVATE_KEY")
	certificate := os.Getenv("GRPC_ADAPTER_CERTIFICATE")
	if credential != "" {
		so, err := getServerTLSOption(credential, privateKey, certificate)
		if err != nil {
			return nil, err
		}
		s.server = grpc.NewServer(so)
	} else {
		s.server = grpc.NewServer()
	}
	metric.RegisterHandleMetricServiceServer(s.server, s)
	return s, nil
}
