open System
open System.IO
open System.Net.Sockets

let serverAddr = "localhost"
let port = 3000

let main() =
    try
        let client = new TcpClient(serverAddr, port)
        let stream = client.GetStream()
        let reader = new StreamReader(stream)
        let writer = new StreamWriter(stream)
        writer.AutoFlush <- true

        let stopReading = ref false

        let readServerResponses () =
            try
                while not !stopReading do
                    let response = reader.ReadLine()
                    if response = null then
                        Console.WriteLine("Connection to the server has been closed.")
                        stopReading := true
                    else
                        Console.WriteLine("Received: {0}", response)
            with
            | :? IOException ->
                Console.WriteLine("Connection to the server has been closed.")
                stopReading := true

        let responseThread = new System.Threading.Thread(readServerResponses)
        responseThread.Start()

        Console.WriteLine("Connected to the server. Type 'exit' to close the connection.")
        let mutable input = ""
        while input <> "exit" do
            input <- Console.ReadLine()
            writer.Write(input)

        Console.WriteLine("Closing the connection.")
        stopReading := true // Signal the reading thread to stop
        client.Close()
    with
    | :? Exception as ex ->
        Console.WriteLine("An error occurred: {0}", ex.Message)

main()

