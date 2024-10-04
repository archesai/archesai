import { Injectable } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { createTransport } from "nodemailer";
import * as Mail from "nodemailer/lib/mailer";

@Injectable()
export class EmailService {
  private nodemailerTransport: Mail;
  constructor(private readonly configService: ConfigService) {
    this.nodemailerTransport = createTransport({
      auth: {
        pass: this.configService.get("EMAIL_PASSWORD"),
        user: this.configService.get("EMAIL_USER"),
      },
      service: this.configService.get("EMAIL_SERVICE"),
    });
  }
  async sendEmailVerification(email: string, link: string) {
    await this.sendMail({
      from: "info@archesai.com",
      html: `<!DOCTYPE html>
<html>
<head>

  <meta charset="utf-8">
  <meta http-equiv="x-ua-compatible" content="ie=edge">
  <title>Email Confirmation</title>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <style type="text/css">
  /**
   * Google webfonts. Recommended to include the .woff version for cross-client compatibility.
   */
  @media screen {
    @font-face {
      font-family: 'Roboto';
      font-style: normal;
      font-weight: 400;
      src: local('Roboto Regular'), local('SourceSansPro-Regular'), url(https://fonts.gstatic.com/s/sourcesanspro/v10/ODelI1aHBYDBqgeIAH2zlBM0YzuT7MdOe03otPbuUS0.woff) format('woff');
    }
    @font-face {
      font-family: 'Roboto';
      font-style: normal;
      font-weight: 700;
      src: local('Roboto Bold'), local('SourceSansPro-Bold'), url(https://fonts.gstatic.com/s/sourcesanspro/v10/toadOcfmlt9b38dHJxOBGFkQc6VGVFSmCnC_l7QZG60.woff) format('woff');
    }
  }
  /**
   * Avoid browser level font resizing.
   * 1. Windows Mobile
   * 2. iOS / OSX
   */
  body,
  table,
  td,
  a {
    -ms-text-size-adjust: 100%; /* 1 */
    -webkit-text-size-adjust: 100%; /* 2 */
  }
  /**
   * Remove extra space added to tables and cells in Outlook.
   */
  table,
  td {
    mso-table-rspace: 0pt;
    mso-table-lspace: 0pt;
  }
  /**
   * Better fluid images in Internet Explorer.
   */
  img {
    -ms-interpolation-mode: bicubic;
  }
  /**
   * Remove blue links for iOS devices.
   */
  a[x-apple-data-detectors] {
    font-family: inherit !important;
    font-size: inherit !important;
    font-weight: inherit !important;
    line-height: inherit !important;
    color: inherit !important;
    text-decoration: none !important;
  }
  /**
   * Fix centering issues in Android 4.4.
   */
  div[style*="margin: 16px 0;"] {
    margin: 0 !important;
  }
  body {
    width: 100% !important;
    height: 100% !important;
    padding: 0 !important;
    margin: 0 !important;
  }
  /**
   * Collapse table borders to avoid space between cells.
   */
  table {
    border-collapse: collapse !important;
  }
  a {
    color: #1a82e2;
  }
  img {
    height: auto;
    line-height: 100%;
    text-decoration: none;
    border: 0;
    outline: none;
  }
  </style>

</head>
<body style="background-color: #e9ecef;">

  <!-- start preheader -->
  <div class="preheader" style="display: none; max-width: 0; max-height: 0; overflow: hidden; font-size: 1px; line-height: 1px; color: #fff; opacity: 0;">
    Welcome to Arches AI! Please confirm your email address.
  </div>
  <!-- end preheader -->

  <!-- start body -->
  <table border="0" cellpadding="0" cellspacing="0" width="100%">

    <!-- start logo -->
    <tr>
      <td align="center" bgcolor="#e9ecef">
        <!--[if (gte mso 9)|(IE)]>
        <table align="center" border="0" cellpadding="0" cellspacing="0" width="600">
        <tr>
        <td align="center" valign="top" width="600">
        <![endif]-->
        <table border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 600px;">
          <tr>
            <td align="center" valign="top" style="padding: 36px 24px;">
              <a href="https://www.archesai.com" target="_blank" style="display: inline-block;">
                <img src="https://www.archesai.com/icon.png" alt="Logo" border="0" width="48" style="display: block; width: 48px; max-width: 48px; min-width: 48px;">
              </a>
            </td>
          </tr>
        </table>
        <!--[if (gte mso 9)|(IE)]>
        </td>
        </tr>
        </table>
        <![endif]-->
      </td>
    </tr>
    <!-- end logo -->

    <!-- start hero -->
    <tr>
      <td align="center" bgcolor="#e9ecef">
        <!--[if (gte mso 9)|(IE)]>
        <table align="center" border="0" cellpadding="0" cellspacing="0" width="600">
        <tr>
        <td align="center" valign="top" width="600">
        <![endif]-->
        <table border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 600px;">
          <tr>
            <td align="left" bgcolor="#ffffff" style="padding: 36px 24px 0; font-family: 'Roboto', Helvetica, Arial; border-top: 3px solid #d4dadf;">
              <h1 style="margin: 0; font-size: 32px; font-weight: 700; letter-spacing: -1px; line-height: 48px;">Welcome to Arches AI,</h1>
            </td>
          </tr>
        </table>
        <!--[if (gte mso 9)|(IE)]>
        </td>
        </tr>
        </table>
        <![endif]-->
      </td>
    </tr>
    <!-- end hero -->

    <!-- start copy block -->
    <tr>
      <td align="center" bgcolor="#e9ecef">
        <!--[if (gte mso 9)|(IE)]>
        <table align="center" border="0" cellpadding="0" cellspacing="0" width="600">
        <tr>
        <td align="center" valign="top" width="600">
        <![endif]-->
        <table border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 600px;">

          <!-- start copy -->
          <tr>
            <td align="left" bgcolor="#ffffff" style="padding: 24px; font-family: 'Roboto', Helvetica, Arial; font-size: 16px; line-height: 24px;">
              <p style="margin: 0;">Thank you for signing up. Please tap the button below to confirm your email address. If you didn't create an account with <a href="https://archesai.com">Arches AI</a>, you can safely delete this email.</p>
            </td>
          </tr>
          <!-- end copy -->

          <!-- start button -->
          <tr>
            <td align="left" bgcolor="#ffffff">
              <table border="0" cellpadding="0" cellspacing="0" width="100%">
                <tr>
                  <td align="center" bgcolor="#ffffff" style="padding: 12px;">
                    <table border="0" cellpadding="0" cellspacing="0">
                      <tr>
                        <td align="center" bgcolor="#2691bf" style="border-radius: 6px;">
                          <a href="${link}" target="_blank" style="display: inline-block; padding: 16px 36px; font-family: 'Roboto', Helvetica, Arial; font-size: 16px; color: #ffffff; text-decoration: none; border-radius: 6px;">Confirm Email Address</a>
                        </td>
                      </tr>
                    </table>
                  </td>
                </tr>
              </table>
            </td>
          </tr>
          <!-- end button -->

          <!-- start copy -->
          <tr>
            <td align="left" bgcolor="#ffffff" style="padding: 24px; font-family: 'Roboto', Helvetica, Arial; font-size: 16px; line-height: 24px; border-bottom: 3px solid #d4dadf">
              <p style="margin: 0;">Cheers,<br> Arches AI Team</p>
            </td>
          </tr>
          <!-- end copy -->

        </table>
        <!--[if (gte mso 9)|(IE)]>
        </td>
        </tr>
        </table>
        <![endif]-->
      </td>
    </tr>
    <!-- end copy block -->

    <!-- start footer -->
    <tr>
      <td align="center" bgcolor="#e9ecef" style="padding: 24px;">
        <!--[if (gte mso 9)|(IE)]>
        <table align="center" border="0" cellpadding="0" cellspacing="0" width="600">
        <tr>
        <td align="center" valign="top" width="600">
        <![endif]-->
        <table border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 600px;">

          <!-- start permission -->
         
          <!-- end permission -->

          <!-- start unsubscribe -->
          <tr>
            <td align="center" bgcolor="#e9ecef" style="padding: 12px 24px; font-family: 'Roboto', Helvetica, Arial; font-size: 14px; line-height: 20px; color: #666;">
              <p style="margin: 0;">This is a no-reply e-mail address. To contact support, please click <a href="mailto:jonathan@archesai.com" target="_blank">here</a></p>  
              <p style="margin: 0;">Arches AI LLC</p>
            </td>
          </tr>
          <!-- end unsubscribe -->

        </table>
        <!--[if (gte mso 9)|(IE)]>
        </td>
        </tr>
        </table>
        <![endif]-->
      </td>
    </tr>
    <!-- end footer -->

  </table>
  <!-- end body -->

</body>
</html>`,
      replyTo: `info@archesai.com`,
      subject: `Arches AI Email Verification - ${new Date().toLocaleString()}`,
      // Write html that looks like this:

      to: email,
    });
  }

  async sendInvite(email: string, link: string) {
    await this.sendMail({
      from: "info@archesai.com",
      html: `<p>Hi there,</p>
      <p>You have been invited to join an organization on Arches AI! Please visit this linke to create your account</p>
      <p>${link}</p>
      <p>Thanks,</p>
      <p>Arches AI Team</p>`,
      replyTo: `info@archesai.com`,
      subject: "Arches AI Invitation",
      to: email,
    });
  }

  async sendMail(options: Mail.Options) {
    return this.nodemailerTransport.sendMail(options);
  }

  async sendPasswordReset(email: string, link: string) {
    await this.sendMail({
      from: "info@archesai.com",
      html: `<!DOCTYPE html>
<html>
<head>

<meta charset="utf-8">
<meta http-equiv="x-ua-compatible" content="ie=edge">
<title>Password Reset</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style type="text/css">
/**
 * Google webfonts. Recommended to include the .woff version for cross-client compatibility.
 */
@media screen {
  @font-face {
    font-family: 'Roboto';
    font-style: normal;
    font-weight: 400;
    src: local('Roboto Regular'), local('SourceSansPro-Regular'), url(https://fonts.gstatic.com/s/sourcesanspro/v10/ODelI1aHBYDBqgeIAH2zlBM0YzuT7MdOe03otPbuUS0.woff) format('woff');
  }
  @font-face {
    font-family: 'Roboto';
    font-style: normal;
    font-weight: 700;
    src: local('Roboto Bold'), local('SourceSansPro-Bold'), url(https://fonts.gstatic.com/s/sourcesanspro/v10/toadOcfmlt9b38dHJxOBGFkQc6VGVFSmCnC_l7QZG60.woff) format('woff');
  }
}
/**
 * Avoid browser level font resizing.
 * 1. Windows Mobile
 * 2. iOS / OSX
 */
body,
table,
td,
a {
  -ms-text-size-adjust: 100%; /* 1 */
  -webkit-text-size-adjust: 100%; /* 2 */
}
/**
 * Remove extra space added to tables and cells in Outlook.
 */
table,
td {
  mso-table-rspace: 0pt;
  mso-table-lspace: 0pt;
}
/**
 * Better fluid images in Internet Explorer.
 */
img {
  -ms-interpolation-mode: bicubic;
}
/**
 * Remove blue links for iOS devices.
 */
a[x-apple-data-detectors] {
  font-family: inherit !important;
  font-size: inherit !important;
  font-weight: inherit !important;
  line-height: inherit !important;
  color: inherit !important;
  text-decoration: none !important;
}
/**
 * Fix centering issues in Android 4.4.
 */
div[style*="margin: 16px 0;"] {
  margin: 0 !important;
}
body {
  width: 100% !important;
  height: 100% !important;
  padding: 0 !important;
  margin: 0 !important;
}
/**
 * Collapse table borders to avoid space between cells.
 */
table {
  border-collapse: collapse !important;
}
a {
  color: #1a82e2;
}
img {
  height: auto;
  line-height: 100%;
  text-decoration: none;
  border: 0;
  outline: none;
}
</style>

</head>
<body style="background-color: #e9ecef;">

<!-- start preheader -->
<div class="preheader" style="display: none; max-width: 0; max-height: 0; overflow: hidden; font-size: 1px; line-height: 1px; color: #fff; opacity: 0;">
  Thank you for using Arches AI! Click this link to reset your password.
</div>
<!-- end preheader -->

<!-- start body -->
<table border="0" cellpadding="0" cellspacing="0" width="100%">

  <!-- start logo -->
  <tr>
    <td align="center" bgcolor="#e9ecef">
      <!--[if (gte mso 9)|(IE)]>
      <table align="center" border="0" cellpadding="0" cellspacing="0" width="600">
      <tr>
      <td align="center" valign="top" width="600">
      <![endif]-->
      <table border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 600px;">
        <tr>
          <td align="center" valign="top" style="padding: 36px 24px;">
            <a href="https://www.archesai.com" target="_blank" style="display: inline-block;">
              <img src="https://www.archesai.com/icon.png" alt="Logo" border="0" width="48" style="display: block; width: 48px; max-width: 48px; min-width: 48px;">
            </a>
          </td>
        </tr>
      </table>
      <!--[if (gte mso 9)|(IE)]>
      </td>
      </tr>
      </table>
      <![endif]-->
    </td>
  </tr>
  <!-- end logo -->

  <!-- start hero -->
  <tr>
    <td align="center" bgcolor="#e9ecef">
      <!--[if (gte mso 9)|(IE)]>
      <table align="center" border="0" cellpadding="0" cellspacing="0" width="600">
      <tr>
      <td align="center" valign="top" width="600">
      <![endif]-->
      <table border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 600px;">
        <tr>
          <td align="left" bgcolor="#ffffff" style="padding: 36px 24px 0; font-family: 'Roboto', Helvetica, Arial; border-top: 3px solid #d4dadf;">
            <h1 style="margin: 0; font-size: 32px; font-weight: 700; letter-spacing: -1px; line-height: 48px;">Welcome to Arches AI,</h1>
          </td>
        </tr>
      </table>
      <!--[if (gte mso 9)|(IE)]>
      </td>
      </tr>
      </table>
      <![endif]-->
    </td>
  </tr>
  <!-- end hero -->

  <!-- start copy block -->
  <tr>
    <td align="center" bgcolor="#e9ecef">
      <!--[if (gte mso 9)|(IE)]>
      <table align="center" border="0" cellpadding="0" cellspacing="0" width="600">
      <tr>
      <td align="center" valign="top" width="600">
      <![endif]-->
      <table border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 600px;">

        <!-- start copy -->
        <tr>
          <td align="left" bgcolor="#ffffff" style="padding: 24px; font-family: 'Roboto', Helvetica, Arial; font-size: 16px; line-height: 24px;">
            <p style="margin: 0;">Please tap the button below to reset your password. If you didn't create an account with <a href="https://archesai.com">Arches AI</a>, you can safely delete this email.</p>
          </td>
        </tr>
        <!-- end copy -->

        <!-- start button -->
        <tr>
          <td align="left" bgcolor="#ffffff">
            <table border="0" cellpadding="0" cellspacing="0" width="100%">
              <tr>
                <td align="center" bgcolor="#ffffff" style="padding: 12px;">
                  <table border="0" cellpadding="0" cellspacing="0">
                    <tr>
                      <td align="center" bgcolor="#2691bf" style="border-radius: 6px;">
                        <a href="${link}" target="_blank" style="display: inline-block; padding: 16px 36px; font-family: 'Roboto', Helvetica, Arial; font-size: 16px; color: #ffffff; text-decoration: none; border-radius: 6px;">Reset Password</a>
                      </td>
                    </tr>
                  </table>
                </td>
              </tr>
            </table>
          </td>
        </tr>
        <!-- end button -->

        <!-- start copy -->
        <tr>
          <td align="left" bgcolor="#ffffff" style="padding: 24px; font-family: 'Roboto', Helvetica, Arial; font-size: 16px; line-height: 24px; border-bottom: 3px solid #d4dadf">
            <p style="margin: 0;">Cheers,<br> Arches AI Team</p>
          </td>
        </tr>
        <!-- end copy -->

      </table>
      <!--[if (gte mso 9)|(IE)]>
      </td>
      </tr>
      </table>
      <![endif]-->
    </td>
  </tr>
  <!-- end copy block -->

  <!-- start footer -->
  <tr>
    <td align="center" bgcolor="#e9ecef" style="padding: 24px;">
      <!--[if (gte mso 9)|(IE)]>
      <table align="center" border="0" cellpadding="0" cellspacing="0" width="600">
      <tr>
      <td align="center" valign="top" width="600">
      <![endif]-->
      <table border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 600px;">

        <!-- start permission -->
       
        <!-- end permission -->

        <!-- start unsubscribe -->
        <tr>
          <td align="center" bgcolor="#e9ecef" style="padding: 12px 24px; font-family: 'Roboto', Helvetica, Arial; font-size: 14px; line-height: 20px; color: #666;">
            <p style="margin: 0;">This is a no-reply e-mail address. To contact support, please click <a href="mailto:jonathan@archesai.com" target="_blank">here</a></p>  
            <p style="margin: 0;">Arches AI LLC</p>
          </td>
        </tr>
        <!-- end unsubscribe -->

      </table>
      <!--[if (gte mso 9)|(IE)]>
      </td>
      </tr>
      </table>
      <![endif]-->
    </td>
  </tr>
  <!-- end footer -->

</table>
<!-- end body -->

</body>
</html>`,
      replyTo: `info@archesai.com`,
      subject: `Arches AI Password Reset - ${new Date().toLocaleString()}`,
      // Write html that looks like this:

      to: email,
    });
  }
}
