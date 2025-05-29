import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";

interface ProcessingOverlayProps {
  isProcessing: boolean;
  setIsProcessing: React.Dispatch<React.SetStateAction<boolean>>;
}

const ProcessingOverlay: React.FC<ProcessingOverlayProps> = ({
  isProcessing,
  setIsProcessing,
}) => {
  const [seconds, setSeconds] = useState<number>(20); // Countdown in seconds
  const navigate = useNavigate();

  useEffect(() => {
    let timerId: NodeJS.Timeout;
    if (isProcessing && seconds > 0) {
      timerId = setInterval(() => {
        setSeconds((prev) => prev - 1);
      }, 1000);
    } else if (seconds === 0) {
      // Redirect to home after countdown ends
      navigate("/dashboard");
      setIsProcessing(false); // Hide overlay after redirection
    }

    return () => clearInterval(timerId); // Cleanup interval on component unmount
  }, [isProcessing, seconds, navigate, setIsProcessing]);

  if (!isProcessing) return null; // Don't render the overlay if not processing

  return (
    <div className="fixed inset-0 bg-gray-900 bg-opacity-50 flex justify-center items-center z-50">
      <div className="bg-white p-8 rounded-lg shadow-lg text-center max-w-sm w-full">
        <h2 className="text-xl font-semibold text-gray-700">Processing...</h2>
        <p className="text-lg text-gray-500 mt-4">
          Please wait, your request is being processed.
        </p>
        <div className="mt-6">
          <span className="text-2xl font-bold text-orange-1">
            {seconds} {seconds <= 1 ? "second" : "seconds"} remaining
          </span>
        </div>
      </div>
    </div>
  );
};

export default ProcessingOverlay;
