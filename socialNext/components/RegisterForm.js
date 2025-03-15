import { useState } from 'react';

const RegisterForm = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [username, setUsername] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (event) => {
    event.preventDefault();

    if (password !== confirmPassword) {
      setError('Passwords do not match!');
      return;
    }

    setError('');
    setLoading(true);

    try {
      const formData = new FormData();
      formData.append('username', username);
      formData.append('email', email);
      formData.append('password', password);

      const response = await fetch('http://localhost:4848', {
        method: 'POST',
        // body: JSON.stringify({
        //     username: 'user123',
        //     password: 'password123'
        // }),
        body: JSON.stringify({ username, password }),
        // body: new URLSearchParams(formData).toString(),
        headers: {
            'Content-Type': 'application/json',
        },
    });

      const data = await response.json();

      if (data.error === 'Validation error' && data.fields) {
        const errorMessages = [];
        // Object.entries(data.fields).forEach(([field, message]) => {
        //   errorMessages.push(message); 
        // });
        setError(errorMessages.join(', ')); 
      } else if (data.error) {
        setError(data.error);
      } else {
        alert(data.message);
        // Reset the form
        setEmail('');
        setPassword('');
        setConfirmPassword('');
        setUsername('');
      }
    } catch (err) {
        console.error('Error:', err); 
        alert('An error occurred.');
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

        <form id="registerForm" className="auth-form" onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="registerUsername">Username</label>
            <input
              type="text"
              id="registerUsername"
              name="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
            {error && <small style={{ color: 'red' }}>{error}</small>}
          </div>
          <div className="form-group">
            <label htmlFor="registerEmail">Email</label>
            <input
              type="email"
              id="registerEmail"
              name="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
            {error && <small style={{ color: 'red' }}>{error}</small>}
          </div>
          <div className="form-group">
            <label htmlFor="registerPassword">Password</label>
            <input
              type="password"
              id="registerPassword"
              name="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
            {error && <small style={{ color: 'red' }}>{error}</small>}
          </div>
          <div className="form-group">
            <label htmlFor="confirmPassword">Confirm Password</label>
            <input
              type="password"
              id="confirmPassword"
              name="confirmPassword"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              required
            />
            {error && <small id="passwordError" style={{ color: 'red' }}>{error}</small>}
          </div>
          <button type="submit" disabled={loading}>
            {loading ? 'Registering...' : 'Register'}
          </button>
        </form>
      </div>
    </div>
  );
};

export default RegisterForm;
