package main

// CloudUploader defines an interface for a type that will be used to upload telemetry to the cloud
type CloudUploader interface {
	SetServer(srv *Server)
	Run()
}
