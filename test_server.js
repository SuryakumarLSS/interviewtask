const express = require('express');
const app = express();
const port = 5001;

app.get('/', (req, res) => res.send('Hello World!'));

const server = app.listen(port, () => {
    console.log(`Test server running on port ${port}`);
});

process.on('exit', (code) => {
    console.log(`Test process exiting with code: ${code}`);
});
