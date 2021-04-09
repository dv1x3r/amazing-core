using System;

namespace Amazing_Tester
{
    class Program
    {
        static void Main(string[] args)
        {
            var session = new GSFSession(null, 0, new GSFBidirectionalSocketTransport(), new GSFBitProtocolCodec(new AWMessageFactory()));
            session.Open("127.0.0.1", 8182);
            //session.SendMessage(ServiceClass.UserServer, 566, new GSFGetClientVersionInfoSvc.GSFRequest("AmazingWorld"), null, null);
            //session.SendMessage(ServiceClass.UserServer, 15, new GSFLoginSvc.GSFRequest("abc", "abc", 1234, new GSFOID(293578400717707620L), "Token", GetEnvironmentData(), null, 0, null), null, null);
            // Create a breakpoint in the GSFSession.cs at line 336
            Console.WriteLine("Hello World!");
        }

        private static GSFClientEnvironmentData GetEnvironmentData()
        {
            GSFClientEnvironmentData gSFClientEnvironmentData = new GSFClientEnvironmentData();
            gSFClientEnvironmentData.unityVersion = "5.3.6p8";
            gSFClientEnvironmentData.userAgent = "main.standalone.133852";
            gSFClientEnvironmentData.screenResolution = "1920x1080";
            gSFClientEnvironmentData.machineOs = "Windows 10  (10.0.0) 64bit;20QGS39V00 (LENOVO);Intel(R) UHD Graphics 620;2111;50;8;16125;web:False";
            gSFClientEnvironmentData.userTime = DateTime.Now;
            gSFClientEnvironmentData.utcOffsetInMinutes = 30;
            gSFClientEnvironmentData.ipAddress = "192.168.0.100";
            return gSFClientEnvironmentData;
        }

    }
}
