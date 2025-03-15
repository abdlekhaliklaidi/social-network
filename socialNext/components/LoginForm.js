import { useState } from 'react';

const LoginForm = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (event) => {
    event.preventDefault();

    if (!email || !password) {
      setError('Email and password are required.');
      return;
    }

    setError('');
    setLoading(true);

    try {
      const formData = new FormData();
      formData.append('email', email);
      formData.append('password', password);

      const response = await fetch('http://localhost:4848/login', {
        method: 'POST',
        // body: new URLSearchParams(formData).toString(),
        body: JSON.stringify({ username, password }),
        headers: {
          'Content-Type': 'application/json',
         },
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'An error occurred');
      }

      // Reset the form and redirect
      setEmail('');
      setPassword('');
      setTimeout(() => {
        window.location.href = '/';
      }, 500);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-popup">
      <div className="auth-container">
        <button className="close-popup" id="closePopup">&times;</button>

        <div className="auth-tabs">
          <button className="active" data-form="login">Login</button>
          <button data-form="register">Register</button>
        </div>

        <form id="loginForm" className="auth-form" onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="loginEmail">Email</label>
            <input
              type="email"
              id="loginEmail"
              name="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
            {error && <small style={{ color: 'red' }}>{error}</small>}
          </div>
          <div className="form-group">
            <label htmlFor="loginPassword">Password</label>
            <input
              type="password"
              id="loginPassword"
              name="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
            {error && <small style={{ color: 'red' }}>{error}</small>}
          </div>
          <button type="submit" disabled={loading}>
            {loading ? 'Logging in...' : 'Login'}
          </button>
        </form>
      </div>
    </div>
  );
};

export default LoginForm;
