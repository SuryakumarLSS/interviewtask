require('dotenv').config();
console.log('Dotenv loaded');
setInterval(() => { console.log('Tick'); }, 1000);
process.on('SIGINT', () => console.log('SIGINT received'));
