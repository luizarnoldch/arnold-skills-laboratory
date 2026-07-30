const express = require('express');
const app = express();

app.use(express.json());

app.post('/api/login', (req, res) => {
  const { username, password } = req.body;
  
  if (!username || !password) {
    return res.status(400).json({ error: 'Username and password are required' });
  }
  
  if (username === 'admin' && password === 'password') {
    return res.json({ token: 'mock-jwt-token', user: { username } });
  }
  
  res.status(401).json({ error: 'Invalid credentials' });
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});

module.exports = app;