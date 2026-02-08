const sqlite3 = require('sqlite3').verbose();
const db = new sqlite3.Database('./rbac.db');

db.serialize(() => {
    // Add status column with default 'Active' (or 'Pending' for new ones)
    // We'll update existing pending invites to 'Pending'
    db.run("ALTER TABLE users ADD COLUMN status TEXT DEFAULT 'Active'", (err) => {
        if (err && !err.message.includes('duplicate column')) console.error("Error adding status:", err.message);
        else console.log("Added status column");

        // Mark users with invitation_token as 'Pending'
        db.run("UPDATE users SET status = 'Pending' WHERE invitation_token IS NOT NULL", (err) => {
            if (err) console.error("Error updating pending status:", err.message);
            else console.log("Updated pending users status");
        });
    });
});

db.close();
