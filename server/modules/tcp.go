package modules


import  (
	"net"
	"log"
)

func Tcp() {

	listener, err := net.Listen("tcp", ":8081")

	if err != nil {
		log.Println("TCP SERVER   FAILED PORT  8081",err)
	}


	defer listener.Close()

	log.Println("TCP server running on port 8081")


	for {


		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}

		go TcpConn(conn) 
		
	}

}

func TcpConn(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("Hello from Go TCP server!\n"))
}
