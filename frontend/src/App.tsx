import React, { useEffect, useState } from 'react';

function App() {
  const [message, setMessage] = useState('');

  useEffect(() => {
    console.log("hello")
    fetch('http://localhost:8080/test.html')
        .then(response => response.text())
        .then(data => setMessage(data))
        .catch(err => console.error(err));
  }, []);

  return (
      <div>
        <h1>React App</h1>
        <p>{message}</p>
      </div>
  );
}

export default App;