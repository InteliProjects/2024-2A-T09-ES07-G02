import { useState, useRef, useEffect } from "react";
import "./index.css";
import logo from "./assets/logo.png";
import { Mic, Settings, Sun, Moon } from "lucide-react";
import LrrContainer from "./components/ui/LrrContainer";
import { RecentLrrsPopover } from "./components/ui/RecentLrrs";
import TagsPopup from "./components/ui/TagsPopup";
import { FilterPopover } from "./components/ui/Filter";

// Type definitions for the SpeechRecognition API
type SpeechRecognition = typeof window.SpeechRecognition | typeof window.webkitSpeechRecognition;
declare global {
  interface Window {
    webkitSpeechRecognition: any;
    SpeechRecognition: any;
  }
}

function App() {
  // State variables for recording, search, and dark mode functionalities
  const [isRecording, setIsRecording] = useState(false);
  const [mediaRecorder, setMediaRecorder] = useState<MediaRecorder | null>(null);
  const [audioChunks, setAudioChunks] = useState<Blob[]>([]);
  const [transcript, setTranscript] = useState("");
  const [searchTerm, setSearchTerm] = useState("");
  const [searchResults, setSearchResults] = useState([]);
  const [showTagsPopup, setShowTagsPopup] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [isDarkMode, setIsDarkMode] = useState(false);

  // Refs for audio and speech recognition
  const audioBlobRef = useRef<Blob | null>(null);
  const recognitionRef = useRef<SpeechRecognition | null>(null);
  const searchInputRef = useRef<HTMLInputElement | null>(null);

  // Toggle dark mode
  const toggleDarkMode = () => {
    setIsDarkMode((prevMode) => !prevMode);
    document.documentElement.classList.toggle("dark", !isDarkMode); // Toggles the "dark" class on the HTML element
  };

  // useEffect for mediaRecorder state changes
  useEffect(() => {
    if (!mediaRecorder) return;

    mediaRecorder.ondataavailable = (event) => {
      if (event.data.size > 0) {
        setAudioChunks((prev) => [...prev, event.data]);
      }
    };

    mediaRecorder.onstop = async () => {
      if (audioChunks.length > 0) {
        const audioBlob = new Blob(audioChunks, {
          type: "audio/webm;codecs=opus",
        });
        audioBlobRef.current = audioBlob;
        sendToBackend(audioBlob, transcript); // Send the recorded audio and transcript to the backend
        setAudioChunks([]);
      } else {
        console.error("No audio data was captured.");
      }
    };
  }, [mediaRecorder, audioChunks, transcript]);

  // useEffect to handle keyboard shortcuts for search input
  useEffect(() => {
    const handleKeyDown = (event: globalThis.KeyboardEvent) => {
      if (event.ctrlKey && event.key === "k") {
        event.preventDefault();
        searchInputRef.current?.focus();
      } else if (event.key === "Escape") {
        searchInputRef.current?.blur();
      }
    };

    window.addEventListener("keydown", handleKeyDown);

    return () => {
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, []);

  // Start speech recognition using the browser's SpeechRecognition API
  const startSpeechRecognition = () => {
    const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition;
    if (!SpeechRecognition) {
      console.error("Speech Recognition API not supported in this browser.");
      return;
    }
    const recognition = new SpeechRecognition();
    recognition.lang = "pt-BR"; // Set language to Brazilian Portuguese
    recognition.interimResults = true;

    recognition.onresult = (event: any) => {
      const lastResultIndex = event.results.length - 1;
      const transcriptText = event.results[lastResultIndex][0].transcript;
      setTranscript(transcriptText); // Update transcript state with the recognized text
    };

    recognition.onend = () => {
      recognition.stop();
    };

    recognition.start();
    recognitionRef.current = recognition;
  };

  // Stop speech recognition
  const stopSpeechRecognition = () => {
    if (recognitionRef.current) {
      recognitionRef.current.stop();
      setSearchTerm((prev) => prev + transcript); 
    }
  };

  // Handles microphone click to start/stop recording and speech recognition
  const handleMicClick = async () => {
    if (!isRecording) {
      try {
        const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
  
        const recorder = new MediaRecorder(stream, { mimeType: "audio/webm;codecs=opus" });
        setMediaRecorder(recorder);
        recorder.start();
        startSpeechRecognition();
        setIsRecording(true);
      } catch (error) {
        console.error("Error accessing the microphone:", error);
      }
    } else {
      mediaRecorder?.stop();
      mediaRecorder?.stream.getTracks().forEach((track) => track.stop());
      stopSpeechRecognition();
      setIsRecording(false);
      setTranscript("");  
    }
  };
  

  // Function to send recorded audio and transcript to the backend
   const sendToBackend = async (audioBlob: Blob, transcript: string) => {
    const formData = new FormData();
    formData.append("audio", audioBlob, "recording.webm");
    formData.append("transcript", transcript);

    try {
      const responseGo = await fetch("http://localhost:7070/upload", {
        method: "POST",
        body: formData,
      });

      if (!responseGo.ok) {
        console.error("Error sending data to the Go backend.");
      }

      const responseCore = await fetch("http://localhost:8080/api/Search", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ transcript }),
      });

      if (!responseCore.ok) {
        console.error("Error sending transcript to C# core.");
      }
    } catch (error) {
      console.error("Request error:", error);
    }
  };

  // Function to handle search submission
  const handleSearchSubmit = async () => {
    if (searchTerm.trim() === "") return;

    setIsLoading(true);
    try {
      const response = await fetch("http://localhost:8080/api/Search", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ text: searchTerm }),
      });

      if (response.ok) {
        const data = await response.json();
        setSearchResults(data);
        console.log("Search results:", data);
      } else {
        console.error("Error in search request:", response.statusText);
      }
    } catch (error) {
      console.error("Error in search request:", error);
    } finally {
      setIsLoading(false);
    }
  };

  // Function to handle the "Enter" key press on search input
  const handleKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === "Enter") {
      handleSearchSubmit();
    }
  };

  // Frontend components and design
  return (
    <div className={`min-h-screen ${isDarkMode ? "bg-gray-900 text-white" : "bg-gray-100 text-black"}`}>
      <header className="flex flex-wrap items-center justify-between p-4 border-b-2 border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900">
        <div className="flex items-center space-x-4">
          <img src={logo} alt="Logo" className="h-16 w-16 sm:h-20 sm:w-20 object-contain" />
          <h1 className="text-xl sm:text-2xl font-bold">LawHunter</h1>
        </div>

        {/* Centers on desktop and remains responsive for mobile */}
        <div className="w-full sm:w-auto mt-4 sm:mt-0 flex justify-between items-center sm:justify-center space-x-4">
          <div className="flex-grow sm:flex-grow-0 relative flex items-center w-full sm:w-96">
            <button className="absolute right-3" onClick={handleMicClick}>
              <Mic className={`h-5 w-5 ${isRecording ? 'text-red-500 animate-pulse' : 'text-black dark:text-black'}`} aria-label="Mic"/>
            </button>
            <input
              ref={searchInputRef}
              type="text"
              placeholder="Digite sua pesquisa..."
              value={isRecording ? transcript : searchTerm}
              onChange={(e) => {
                if (!isRecording) setSearchTerm(e.target.value);  // Permite a edição manual somente quando não está gravando
              }}
              onKeyDown={handleKeyDown}
              className="w-full pl-4 pr-16 py-3 border dark:text-black border-gray-300 dark:border-gray-700 rounded-full focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <button className="bg-blue-500 hover:bg-blue-700 text-white font-semibold py-2 px-6 rounded-full shadow-md transition-all duration-300 sm:ml-4" onClick={handleSearchSubmit}>
            Pesquisar
          </button>
        </div>

        {/* Align icons on mobile to the right */}
        <div className="flex items-center space-x-3 mt-4 sm:mt-0">
          <RecentLrrsPopover />
          <FilterPopover />
          <button className="relative p-2 rounded-full text-black dark:text-white" onClick={() => setShowTagsPopup(!showTagsPopup)} aria-label="Settings">
            <Settings className="h-5 w-5 text-black dark:text-white" />
          </button>

          {/* Toggle Dark Mode */}
          <button onClick={toggleDarkMode} className="p-2 rounded-full text-black dark:text-white transition-colors duration-300">
            {isDarkMode ? <Sun className="h-5 w-5 text-white" /> : <Moon className="h-5 w-5 text-black" />}
          </button>
        </div>
      </header>

      <main className="flex-grow p-4">
        <LrrContainer searchTerm={searchTerm} searchResults={searchResults} isLoading={isLoading} />
      </main>

      {showTagsPopup && <TagsPopup onClose={() => setShowTagsPopup(false)} />}
    </div>
  );
}

export default App;
