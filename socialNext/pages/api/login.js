// pages/api/login.js
import { serialize } from 'cookie';

export default async function handler(req, res) {
  if (req.method === 'POST') {
    const { email, password } = req.body;

    if (!email || !password) {
      return res.status(400).json({ message: 'Email et mode passe important!' });
    }

    try {
      const response = await fetch('http://localhost:4848/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });

      const data = await response.json();

      if (response.ok) {
        res.setHeader('Set-Cookie', serialize('auth_token', data.token, { path: '/' }));
        return res.status(200).json({ message: 'login succete' });
      } else {
        return res.status(401).json({ message: 'Email et mode passe inccorect' });
      }
    } catch (error) {
      console.log(error);    
      return res.status(500).json({ message: 'error login' });
    }
  } else {
    return res.status(405).json({ message: 'methode failed' });
  }
}
