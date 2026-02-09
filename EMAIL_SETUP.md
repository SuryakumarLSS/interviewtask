# Email Setup Guide

## ✅ Email Sending Configured

Your application is now configured to send **real invitation emails** via Gmail SMTP.

## Current Configuration

The following settings are configured in your `.env` file:

```env
SMTP_USER=jeyaroshini2001@gmail.com
SMTP_PASSWORD=zwdf lcok ljks cjsl
FROM_EMAIL=jeyaroshini2001@gmail.com
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
```

## How It Works

1. When you click **"Send Invitation"** in the Admin Console, the system will:
   - Create a user account with "Pending" status
   - Generate a unique invitation token
   - Send an HTML email to the recipient's email address
   - Log success/failure in the server console

2. The recipient will receive a professional email with:
   - A clickable "Set Your Password" button
   - A copy of the invitation link
   - Professional styling

## Testing

1. Navigate to **Admin Console** → **Invite User**
2. Enter an email address and select a role
3. Click **Send Invitation**
4. Check the recipient's inbox (or use your own email for testing)

## Troubleshooting

### Email Not Received?

**Check Gmail App Password:**
- Gmail requires an "App Password" if you have 2-Factor Authentication enabled
- Your current password looks like an App Password (format: `xxxx xxxx xxxx xxxx`)
- If emails fail, you may need to regenerate your App Password:
  1. Go to https://myaccount.google.com/apppasswords
  2. Generate a new App Password for "Mail"
  3. Update `SMTP_PASSWORD` in `.env` with the new password
  4. Restart the server

**Check Server Logs:**
- Look for `✅ Invitation email sent to: ...` (success)
- Or `❌ Email send failed: ...` (error with details)

**Gmail Security:**
- Ensure "Less secure app access" is OFF (you should use App Passwords instead)
- Check that your Google account doesn't have unusual activity blocks

### Alternative Email Providers

If you prefer to use a different email service:

**SendGrid (Free tier: 100 emails/day):**
```env
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASSWORD=your_sendgrid_api_key
FROM_EMAIL=your_verified_sender@domain.com
```

**Mailgun:**
```env
SMTP_HOST=smtp.mailgun.org
SMTP_PORT=587
SMTP_USER=your_mailgun_smtp_username
SMTP_PASSWORD=your_mailgun_smtp_password
FROM_EMAIL=your_verified_sender@domain.com
```

**Outlook/Hotmail:**
```env
SMTP_HOST=smtp-mail.outlook.com
SMTP_PORT=587
SMTP_USER=your_email@outlook.com
SMTP_PASSWORD=your_password
FROM_EMAIL=your_email@outlook.com
```

## Production Recommendations

For **production use**, consider:

1. **Use a Dedicated Email Service**:
   - SendGrid, Mailgun, AWS SES, or Postmark
   - Better deliverability and analytics
   - No Gmail sending limits

2. **Environment Variables**:
   - Keep `.env` out of version control (add to `.gitignore`)
   - Use environment-specific configurations

3. **Custom Domain**:
   - Use a custom domain for `FROM_EMAIL` (e.g., `noreply@yourdomain.com`)
   - Improves trust and deliverability

4. **Email Templates**:
   - Store HTML templates separately
   - Support multiple languages
   - Add company branding

## Ready to Test! 🚀

Your email system is now live. Try sending an invitation to see it in action!
