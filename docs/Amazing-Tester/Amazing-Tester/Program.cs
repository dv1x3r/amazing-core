using System;

namespace Amazing_Tester
{
    class Program
    {
        static void Main(string[] args)
        {
            var session = new GSFSession(null, 0, new GSFBidirectionalSocketTransport(), new GSFBitProtocolCodec(new AWMessageFactory()));
            session.Open("127.0.0.1", 8182);
            session.SendMessage(ServiceClass.UserServer, 566, new GSFGetClientVersionInfoSvc.GSFRequest("AmazingWorld"), null, null);
            // Create a breakpoint in the GSFSession.cs at line 336
            Console.WriteLine("Hello World!");
        }
    }
}
