import { useEffect, useRef, useState } from "react";
import "./App.css";
import { Greet, ProcessMessage, ToggleVisibility } from "../wailsjs/go/main/App";
import { EventsOn } from "../wailsjs/runtime/runtime.js";

function App() {
    const [resultText, setResultText] = useState("");
    const [name, setName] = useState("");
    const inputRef = useRef(null); // Create a reference for the input field

    // Update the name state
    const updateName = (e) => setName(e.target.value);

    // Process the message when Enter is pressed
    const handleKeyDown = (e) => {
        if (e.key === "Enter") {
            processMessage();
        }

        if (e.key === "Escape") {
            e.preventDefault(); // <-- No more bonk sounds
            ToggleVisibility();
            setName(""); // <-- clears input field
        }
    };

    // Clears the input and can interact with backend if needed
    const processMessage = () => {
        if (!name.trim()) {
            setResultText("⚠️ Input cannot be empty."); // Show a warning for empty input
            return;
        }

        ProcessMessage(name)
            .then(() => {
                setName(""); // Clear input field
            })
            .catch(() => {
                setResultText("❌ An error occurred while processing the message.");
            });
    };

    // Listen for app focus events and focus the input field
    useEffect(() => {
        const focusInput = () => {
            if (inputRef.current) {
                inputRef.current.focus();
                setResultText("");
            }
        };

        setTimeout(() => {
            if (document.hasFocus()) {
                focusInput();
            }
        }, 50);

        EventsOn("wails:focus", focusInput);
    }, []);

    // Listen for backend error events
    useEffect(() => {
        EventsOn("Backend:ErrorEvent", (errorMessage) => {
            setResultText(`❌ Error: ${errorMessage}`);
        });
    }, []);

    return (
        <div className="spotlight-container">
            {/* Spotlight-style input bar */}
            <div className="spotlight-box">
                <input
                    ref={inputRef} // Attach the reference to the input field
                    className="spotlight-input"
                    type="text"
                    placeholder="Type your task..."
                    value={name}
                    onKeyDown={handleKeyDown}
                    onChange={updateName}
                    autoComplete="off"
                />
            </div>

            {/* Render result text below input bar */}
            {resultText && <div className="spotlight-results">{resultText}</div>}
        </div>
    );
}

export default App;