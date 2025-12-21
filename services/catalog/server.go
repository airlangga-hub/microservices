package catalog

type Server struct {
	
	Svc Service
}

func ListenGrpc(service Service, port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("ERROR: account server ListenGrpc: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterAccountServiceServer(s, &Server{Svc: service})

	return s.Serve(lis)
}