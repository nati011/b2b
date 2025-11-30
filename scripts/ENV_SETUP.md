# Environment Variables Setup

## Google OAuth Configuration

The following Google OAuth credentials have been configured:

- **GOOGLE_CLIENT_ID**: `286726142980-uvab0hji7kpsfi8r0c67b7kdh18hs48e.apps.googleusercontent.com`
- **GOOGLE_CLIENT_SECRET**: `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`

## Configuration Files

### Docker Compose
The credentials are configured in `docker-compose.yaml` with default values. You can override them by:

1. **Using environment variables** (recommended for production):
   ```bash
   export GOOGLE_CLIENT_ID=286726142980-uvab0hji7kpsfi8r0c67b7kdh18hs48e.apps.googleusercontent.com
   export GOOGLE_CLIENT_SECRET=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
   docker-compose up
   ```

2. **Using a .env file** (create `.env` in project root):
   ```env
   GOOGLE_CLIENT_ID=286726142980-uvab0hji7kpsfi8r0c67b7kdh18hs48e.apps.googleusercontent.com
   GOOGLE_CLIENT_SECRET=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
   ```

### Services Configured

1. **Client Service** (`web/client`)
   - Port: 3001
   - Google OAuth enabled for customer-facing app

2. **Admin Service** (`web/admin`)
   - Port: 3000
   - Google OAuth enabled for admin panel

## Google Cloud Console Setup

Make sure the following redirect URIs are configured in Google Cloud Console:

### Client App (Port 3001)
- `http://localhost:3001/api/auth/callback/google`
- `https://your-production-domain.com/api/auth/callback/google`

### Admin App (Port 3000)
- `http://localhost:3000/api/auth/callback/google`
- `https://your-admin-domain.com/api/auth/callback/google`

## Testing

After updating the configuration:

1. Restart the services:
   ```bash
   docker-compose down
   docker-compose up --build
   ```

2. Test Google Sign-In:
   - Navigate to `http://localhost:3001/login` (client)
   - Navigate to `http://localhost:3000/auth/signin` (admin)
   - Click "Sign in with Google"
   - Verify the OAuth flow works correctly

## Security Notes

⚠️ **Important**: 
- Never commit `.env` files with real credentials to version control
- Use environment variables or secrets management in production
- Rotate credentials if they are exposed
- The default values in docker-compose.yaml are for development only


