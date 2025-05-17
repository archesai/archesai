export function getEmailChangeConfirmationHtml(
  changeEmailLink: string,
  currentEmail: string
): string {
  return `
    <div style="font-family: Arial, sans-serif; line-height: 1.6;">
      <h2>Hello!</h2>
      <p>You requested to change your email address from ${currentEmail}. Please confirm the change by clicking the link below:</p>
      <a href="${changeEmailLink}" style="display: inline-block; padding: 10px 20px; background-color: #2563eb; color: white; text-decoration: none; border-radius: 5px;">
        Confirm Email Change
      </a>
      <p>If you did not request this change, please ignore this email.</p>
      <p>Best regards,<br/>Arches AI</p>
    </div>
  `
}

export function getEmailVerificationHtml(verificationLink: string): string {
  return `
    <div style="font-family: Arial, sans-serif; line-height: 1.6;">
      <h2>Hello!</h2>
      <p>Thank you for registering. Please verify your email address by clicking the link below:</p>
      <a href="${verificationLink}" style="display: inline-block; padding: 10px 20px; background-color: #2563eb; color: white; text-decoration: none; border-radius: 5px;">
        Verify Email
      </a>
      <p>If you did not create an account, you can safely ignore this email.</p>
      <p>Best regards,<br/>Arches AI</p>
    </div>
  `
}

export function getPasswordResetHtml(resetLink: string): string {
  return `
    <div style="font-family: Arial, sans-serif; line-height: 1.6;">
      <h2>Hello!</h2>
      <p>You have requested a password reset. Click the link below to set a new password:</p>
      <a href="${resetLink}" style="display: inline-block; padding: 10px 20px; background-color: #2563eb; color: white; text-decoration: none; border-radius: 5px;">
        Reset Password
      </a>
      <p>If you did not request this, please ignore this email.</p>
      <p>Best regards,<br/>Arches AI</p>
    </div>
  `
}
