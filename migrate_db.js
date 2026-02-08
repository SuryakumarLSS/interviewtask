const sqlite3 = require('sqlite3').verbose();
const db = new sqlite3.Database('./rbac.db');

db.serialize(() => {
    // Add email column
    db.run("ALTER TABLE users ADD COLUMN email TEXT", (err) => {
        if (err && !err.message.includes('duplicate column')) console.error("Error adding email:", err.message);
        else console.log("Added email column");
    });

    // Add invitation_token column
    db.run("ALTER TABLE users ADD COLUMN invitation_token TEXT", (err) => {
        if (err && !err.message.includes('duplicate column')) console.error("Error adding invitation_token:", err.message);
        else console.log("Added invitation_token column");
    });
});

db.close();
