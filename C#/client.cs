using System;
using System.Diagnostics;
using System.Net.Http;
using System.Threading.Tasks;

class Program
{
    static async Task Main(string[] args)
    {
        string mp3Url = "http://localhost:8080/mp3/your_mp3_file.mp3"; // Replace with the desired mp3 file path
        using HttpClient httpClient = new HttpClient();

        try
        {
            var response = await httpClient.GetStreamAsync(mp3Url);
            using (var player = new Process())
            {
                player.StartInfo.FileName = "mpg123"; // You need to have mpg123 installed on your system
                player.StartInfo.Arguments = "--quiet -"; // Read from stdin
                player.StartInfo.UseShellExecute = false;
                player.StartInfo.RedirectStandardInput = true;
                player.StartInfo.RedirectStandardOutput = false;
                player.StartInfo.RedirectStandardError = false;

                player.Start();
                await response.CopyToAsync(player.StandardInput.BaseStream);
                player.WaitForExit();
            }
        }
        catch (Exception e)
        {
            Console.WriteLine($"Error: {e.Message}");
        }
    }
}