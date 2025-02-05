import React, { useEffect, useState } from 'react';
import axios from 'axios';

const HealthStatus = () => {
  const [status, setStatus] = useState(null);
  const [time, setTime] = useState(null);
  const [error, setError] = useState(null);

  useEffect(() => {
    axios.get('/api/health')
      .then(response => {
        setStatus(response.data.status);
        setTime(new Date(response.data.time * 1000).toLocaleString());
      })
      .catch(error => {
        setError(error.message);
      });
  }, []);

  if (error) {
    return <div>Error: {error}</div>;
  }

  if (!status) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h1>Health Status</h1>
      <p>Status: {status}</p>
      <p>Time: {time}</p>
    </div>
  );
};

export default HealthStatus;