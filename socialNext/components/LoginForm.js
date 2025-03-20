// import { useState } from 'react';

// const LoginForm = () => {
//   const [email, setEmail] = useState('');
//   const [password, setPassword] = useState('');
//   const [error, setError] = useState('');
//   const [loading, setLoading] = useState(false);

//   const handleSubmit = async (event) => {
//     event.preventDefault();

//     if (!email || !password) {
//       setError('Email and password are required.');
//       return;
//     }

//     setError('');
//     setLoading(true);

//     try {
//       const formData = new FormData();
//       formData.append('email', email);
//       formData.append('password', password);

//       const response = await fetch('http://localhost:4848/login', {
//         method: 'POST',
//         // body: new URLSearchParams(formData).toString(),
//         body: JSON.stringify({ username, password }),
//         headers: {
//           'Content-Type': 'application/json',
//          },
//       });

//       const data = await response.json();

//       if (!response.ok) {
//         throw new Error(data.error || 'An error occurred');
//       }

//       // Reset the form and redirect
//       setEmail('');
//       setPassword('');
//       setTimeout(() => {
//         window.location.href = '/';
//       }, 500);
//     } catch (err) {
//       setError(err.message);
//     } finally {
//       setLoading(false);
//     }
//   };

//   return (
//     <div className="auth-popup">
//       <div className="auth-container">
//         <button className="close-popup" id="closePopup">&times;</button>

//         <div className="auth-tabs">
//           <button className="active" data-form="login">Login</button>
//           <button data-form="register">Register</button>
//         </div>

//         <form id="loginForm" className="auth-form" onSubmit={handleSubmit}>
//           <div className="form-group">
//             <label htmlFor="loginEmail">Email</label>
//             <input
//               type="email"
//               id="loginEmail"
//               name="email"
//               value={email}
//               onChange={(e) => setEmail(e.target.value)}
//               required
//             />
//             {error && <small style={{ color: 'red' }}>{error}</small>}
//           </div>
//           <div className="form-group">
//             <label htmlFor="loginPassword">Password</label>
//             <input
//               type="password"
//               id="loginPassword"
//               name="password"
//               value={password}
//               onChange={(e) => setPassword(e.target.value)}
//               required
//             />
//             {error && <small style={{ color: 'red' }}>{error}</small>}
//           </div>
//           <button type="submit" disabled={loading}>
//             {loading ? 'Logging in...' : 'Login'}
//           </button>
//         </form>
//       </div>
//     </div>
//   );
// };

// export default LoginForm;

import { useState } from 'react';

const LoginForm = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');

    try {
      const response = await fetch('http://localhost:4848/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      });

      if (!response.ok) {
        throw new Error('Login failed');
      }

      const data = await response.json();
      console.log('Login successful:', data);
    } catch (err) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div>
      <h1>Login</h1>
      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="username">Username</label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
        </div>
        <div>
          <label htmlFor="password">Password</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <button type="submit" disabled={isLoading}>
          {isLoading ? 'Logging in...' : 'Login'}
        </button>
      </form>
      {error && <p>{error}</p>}
    </div>
  );
};

export default LoginForm;
