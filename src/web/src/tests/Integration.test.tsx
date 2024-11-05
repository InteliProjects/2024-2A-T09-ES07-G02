describe("Integration Tests", () => {
    it("should send data to the backend and receive a response", async () => {
      const audioBlob = new Blob(); 
      const transcript = "Test transcript";
  
      // Create a FormData object
      const formData = new FormData();
      formData.append("audio", audioBlob, "recording.webm");
      formData.append("transcript", transcript);
  
      // Perform the call to your real endpoint
      const responseGo = await fetch("http://localhost:7070/upload", {
        method: "POST",
        body: formData,
      });
  
      expect(responseGo.ok).toBe(true); // Check if the response is 200 OK
  
      const responseCore = await fetch("http://localhost:8080/api/Search", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ transcript }),
      });
  
      expect(responseCore.ok).toBe(true); // Check if the response is 200 OK
    });
  });
  