-- roles 

-- resources

-- superadmin user

-- create payment partners
SELECT public.create_payment_partner(
    'chapa', 
    'icon', 
    'ACTIVE', 
    'https://api.chapa.co', 
    'CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5'
    'DIGITAL_PAYMENT'
);

SELECT public.create_payment_partner(
    'payOnDelivery', 
    'icon', 
    'ACTIVE', 
    'https://no.baseurl', 
    'no-secret',
    'MANUAL_PAYMENT'
);

-- create email template
-- SELECT public.create_email_template(
--     'Password Reset | Ahoy, Climb aboard!',
--     '<!DOCTYPE html>\n<html>\n<head>\n    <style>\n        body { font-family: Arial, sans-serif; line-height: 1.6; }\n        .container { max-width: 600px; margin: 0 auto; padding: 20px; }\n        .button {\n            display: inline-block;\n            background-color: #007BFF;\n            color: white !important;\n            padding: 10px 20px;\n            text-decoration: none;\n            border-radius: 5px;\n        }\n        .footer { margin-top: 20px; font-size: 12px; color: #666; }\n    </style>\n</head>\n<body>\n    <div class="container">\n        <h2>Password Reset Request</h2>\n        <p>Hello {{.username}},</p>\n        <p>We received a request to reset your password. Click the button below to proceed:</p>\n        <p>\n            <a href="google.com/{{.reset_link}}" class="button">Reset Password</a>\n        </p>\n        <p>If you didnt request this, please ignore this email.</p>\n        <div class="footer">\n            <p>Best regards,<br>The Team</p>\n        </div>\n    </div>\n</body>\n</html>'
-- );