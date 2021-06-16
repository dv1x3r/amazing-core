using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Console_Tester
{
    class Program
    {
        static void Main(string[] args)
        {
            var session = new GSFSession(null, 0, new GSFBidirectionalSocketTransport(), new GSFBitProtocolCodec(new AWMessageFactory()));
            session.Open("127.0.0.1", 8182); // Create a breakpoint in the GSFSession.cs at line 336
            //session.SendMessage(ServiceClass.UserServer, 566, new GSFGetClientVersionInfoSvc.GSFRequest("AmazingWorld"), null, null);
            //session.SendMessage(ServiceClass.UserServer, 213, new GSFGetRequiredExperienceSvc.GSFRequest(1), null, null);
            //session.SendMessage(ServiceClass.UserServer, 103, new GSFGetOutfitItemsSvc.GSFRequest(new GSFOID(), new GSFOID()), null, null);
            session.SendMessage(ServiceClass.UserServer, 154, new GSFGetZonesSvc.GSFRequest(), null, null);
            Console.WriteLine("Hello World!");
        }
    }
}
