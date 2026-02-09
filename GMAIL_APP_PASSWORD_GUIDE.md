# 🔐 Gmail App Password Setup Guide

## Current Error
```
535 5.7.8 Username and Password not accepted
```

This means Gmail is rejecting the App Password. Let's get you a fresh one!

## ✅ Step-by-Step: Generate New Gmail App Password

### Step 1: Enable 2-Factor Authentication (If Not Already)
1. Go to: https://myaccount.google.com/security
2. Under "How you sign in to Google", select **2-Step Verification**
3. If not enabled, click **Get Started** and follow the setup

### Step 2: Generate App Password
1. Go to: https://myaccount.google.com/apppasswords
   - OR Google Account → Security → 2-Step Verification → App passwords
2. You might need to sign in again
3. Under "Select app", choose **Mail**
4. Under "Select device", choose **Other (Custom name)**
5. Enter: `RBAC System` (or any name you want)
6. Click **Generate**
7. **COPY THE 16-CHARACTER PASSWORD** (example: `abcd efgh ijkl mnop`)

### Step 3: Update Your .env File
1. Open `d:\project-new\.env`
2. Update line 2 with the new password **WITHOUT SPACES**:
   ```env
   SMTP_PASSWORD=abcdefghijklmnop
   ```
   Example: If Google shows `abcd efgh ijkl mnop`, you enter `abcdefghijklmnop`

3. Save the file

### Step 4: Restart the Server
```bash
# In your terminal (or I can do this for you)
npm run dev
```

### Step 5: Test
1. Go to Admin Console
2. Send an invitation
3. Check your inbox!

## 🔧 Alternative: Use a Different Email Provider

If you prefer not to use Gmail, here are other options:

### Option 1: Outlook/Hotmail (Easiest)
```env
SMTP_HOST=smtp-mail.outlook.com
SMTP_PORT=587
SMTP_USER=your_email@outlook.com
SMTP_PASSWORD=your_actual_password
FROM_EMAIL=your_email@outlook.com
```
✅ No App Password needed
✅ Just use your regular Outlook password

### Option 2: SendGrid (Best for Production)
1. Sign up at: https://sendgrid.com (Free tier: 100 emails/day)
2. Create an API key
3. Update .env:
   ```env
   SMTP_HOST=smtp.sendgrid.net
   SMTP_PORT=587
   SMTP_USER=apikey
   SMTP_PASSWORD=your_sendgrid_api_key
   FROM_EMAIL=your_verified@email.com
   ```

### Option 3: Mailtrap (Best for Testing)
Perfect for development - catches all emails in a fake inbox!
1. Sign up at: https://mailtrap.io (Free)
2. Get SMTP credentials from your inbox
3. Update .env with provided credentials

## 🐛 Still Having Issues?

### Common Problems:

**Problem**: "Less secure app access"
- **Solution**: Gmail deprecated this. You MUST use App Passwords with 2FA enabled

**Problem**: "Username and Password not accepted"
- **Solution**: Double-check the password has NO SPACES
- **Solution**: Generate a new App Password

**Problem**: "No .env file found"
- **Solution**: Make sure `.env` is in `d:\project-new\` (project root)

## 📝 Current Configuration

Your current `.env` file location: `d:\project-new\.env`

Current settings:
```env
SMTP_USER=jeyaroshini2001@gmail.com
SMTP_PASSWORD=zwdflcokljkscjsl  ← UPDATE THIS
FROM_EMAIL=jeyaroshini2001@gmail.com
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
```

## 🎯 Next Steps

1. **Generate a fresh App Password** using the steps above
2. **Update line 2** in `.env` with the new password (no spaces)
3. **Tell me when done**, and I'll restart the server for you
4. **Test sending an invitation**

The password I currently have (`zwdflcokljkscjsl`) may not be correct. Please generate a fresh one following the steps above!
