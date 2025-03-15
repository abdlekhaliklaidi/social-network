import React, { useState } from 'react';

const AuthForm = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');

  const handleLogin = async () => {
    const response = await fetch('http://localhost:4848/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email,
        password,
      }),
    });

    if (response.ok) {
      const data = await response.json();
      setMessage('Login successful');
      console.log(data);
    } else {
      setMessage('Login failed');
    }
  };

  const handleRegister = async () => {
    const response = await fetch('http://localhost:4848/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email,
        password,
      }),
    });

    if (response.ok) {
      const data = await response.json();
      setMessage('Registration successful');
      console.log(data);
    } else {
      setMessage('Registration failed');
    }
  };

  const handleLogout = async () => {
    const response = await fetch('http://localhost:4848/logout', {
      method: 'POST',
    });

    if (response.ok) {
      const data = await response.json();
      setMessage('Logout successful');
      console.log(data);
    } else {
      setMessage('Logout failed');
    }
  };

  return (
    <div>
      <h2>Auth Form</h2>
      <input
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        placeholder="Email"
      />
      <input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="Password"
      />
      <button onClick={handleLogin}>Login</button>
      <button onClick={handleRegister}>Register</button>
      <button onClick={handleLogout}>Logout</button>
      <p>{message}</p>
    </div>
  );
};

export default AuthForm;
