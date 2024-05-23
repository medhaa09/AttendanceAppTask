import React, { useState, useRef } from 'react';
import Webcam from 'react-webcam';
import axios from 'axios';

const Home = () => {
  const [result, setResult] = useState("");
  const webcamRef = useRef(null);

  const capture = async () => {
    const imageSrc = webcamRef.current.getScreenshot();
    if (imageSrc) {
      const base64Image = imageSrc.split(',')[1]; // Extract base64 part
      try {
        const response = await axios.post('http://localhost:8080/mark-attendance', { image: base64Image }, {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + localStorage.getItem('token')
          }
        });
        setResult(response.data);
      } catch (error) {
        console.error('Error recognizing face:', error);
      }
    } else {
      console.error('Failed to capture image');
    }
  };
  
  return (
    <div>
      <h1>Mark Attendance</h1>
      <Webcam
        audio={false}
        ref={webcamRef}
        screenshotFormat="image/jpeg"
        width={320}
        height={240}
      />
      <button onClick={capture}>Capture</button>
      {result && <pre>{JSON.stringify(result, null, 2)}</pre>}
    </div>
  );
};

export default Home;
