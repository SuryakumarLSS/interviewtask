const express = require('express');
const cors = require('cors');
const app = express();
require('dotenv').config();

app.use(cors());
app.use(express.json());

const db = require('./database');
const authMiddleware = require('./middleware/auth');

app.use('/api/auth', require('./routes/auth'));
app.use('/api/admin', authMiddleware, require('./routes/admin'));
app.use('/api/data', authMiddleware, require('./routes/data'));

const PORT = 5001;
const server = app.listen(PORT, () => {
    console.log(`Server running on port ${PORT}`);
});

// Force keep-alive to ensure server stability
setInterval(() => {
    // Heartbeat
}, 60000);

process.on('exit', (code) => {
    console.log(`Process exiting with code: ${code}`);
});

process.on('SIGINT', () => {
    console.log('Received SIGINT. Shutting down.');
    server.close(() => {
        console.log('Server closed.');
        process.exit(0);
    });
});
