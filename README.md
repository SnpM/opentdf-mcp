# OpenTDF MCP Gateway - End-to-End Encryption Demo

This repository contains an end-to-end demo showcasing OpenTDF's file encryption and decryption capabilities using the [@opentdf/sdk](https://www.npmjs.com/package/@opentdf/sdk) web SDK.

## 🚀 Quick Start

1. **Install dependencies:**
   ```bash
   npm install
   ```

2. **Start the demo server:**
   ```bash
   npm run dev
   ```

3. **Open your browser:**
   Navigate to `http://localhost:8080` to access the demo interface.

## 📋 Features

### ✨ What's Included

- **Web-based Interface**: Clean, responsive UI for easy interaction
- **File Encryption**: Encrypt any file using OpenTDF's NanoTDF format
- **File Decryption**: Decrypt TDF files back to their original content
- **Sample File Generation**: Create test files for demonstration
- **Configuration Management**: Set up your OpenTDF service endpoints
- **Download Support**: Download encrypted and decrypted files

### 🔧 Technical Features

- Uses OpenTDF Web SDK (`@opentdf/sdk`)
- NanoTDF format for efficient encryption
- Stream-based processing for large files
- Browser-compatible implementation
- No server-side dependencies required

## 🛠️ Configuration

The demo requires configuration of your OpenTDF service endpoints:

- **KAS Endpoint**: Your Key Access Service URL
- **OIDC Origin**: Your OpenID Connect provider URL  
- **Client ID**: Your registered client identifier
- **Refresh Token**: (Optional) For authenticated requests

### Default Demo Configuration

The demo comes pre-configured with demo endpoints:
- KAS: `https://demo.kas.opentdf.io`
- OIDC: `https://demo.keycloak.opentdf.io`
- Client ID: `demo-client`

## 📖 How to Use

### 1. Encrypt a File

1. Configure your OpenTDF service endpoints (or use the defaults)
2. Click "Create Sample Text File" or "Create Sample JSON File" to generate test data
3. Use the file picker to select your file
4. Click "Encrypt File" to create an encrypted TDF file
5. Download the encrypted `.ntdf` file

### 2. Decrypt a File

1. Use the file picker to select an encrypted `.tdf` or `.ntdf` file
2. Click "Decrypt File" to decrypt the content
3. View the decrypted content in the result area
4. Download the decrypted file if needed

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Browser   │    │  OpenTDF SDK    │    │ OpenTDF Service │
│                 │────│                 │────│                 │
│  - demo.js      │    │ - Encryption    │    │ - Key Access    │
│  - index.html   │    │ - Decryption    │    │ - OIDC Auth     │
│                 │    │ - Auth Provider │    │ - Policies      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Key Components

- **Frontend (Browser)**: HTML interface with JavaScript demo logic
- **OpenTDF SDK**: Handles encryption/decryption operations
- **OpenTDF Service**: Provides key management and access control

## 🔐 Security Features

- **Zero-Trust Architecture**: Data is encrypted with policy-based access controls
- **Attribute-Based Access Control (ABAC)**: Fine-grained permissions
- **End-to-End Encryption**: Data remains encrypted in transit and at rest
- **Key Management**: Centralized key access service
- **Policy Enforcement**: Server-side policy evaluation

## 📁 File Structure

```
opentdf-mcp/
├── index.html          # Main demo interface
├── demo.js            # Demo implementation logic
├── package.json       # Node.js dependencies
├── README.md          # This documentation
├── DESIGN.md          # Original project design
└── node_modules/      # Dependencies (after npm install)
    └── @opentdf/sdk/  # OpenTDF Web SDK
```

## 🧪 Development

### Available Scripts

- `npm start` - Start HTTP server on port 8080
- `npm run dev` - Start server and open browser automatically

### Adding New Features

The demo is built with vanilla JavaScript modules for simplicity. To extend:

1. **Add new encryption formats**: Modify `demo.js` to support ZTDF or TDF3
2. **Enhance UI**: Update `index.html` with additional interface elements
3. **Add policy support**: Integrate attribute-based policies into encryption
4. **Improve auth**: Implement proper OIDC authentication flow

## 🐛 Troubleshooting

### Common Issues

1. **"Failed to initialize client"**
   - Check your OpenTDF service endpoints
   - Verify network connectivity
   - Ensure CORS is properly configured

2. **"Encryption/Decryption failed"**
   - Verify your authentication credentials
   - Check that the KAS endpoint is accessible
   - Ensure the file format is supported

3. **Browser compatibility**
   - Modern browsers required (ES6+ support)
   - HTTPS required for production deployments
   - File API and Streams API support needed

### Debug Mode

Open browser developer tools (F12) to see detailed error messages and network requests.

## 🤝 Contributing

This is a demo project for the OpenTDF MCP Gateway. Contributions welcome:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📚 Resources

- [OpenTDF Platform](https://github.com/opentdf/platform)
- [OpenTDF Web SDK](https://www.npmjs.com/package/@opentdf/sdk)
- [OpenTDF Documentation](https://docs.opentdf.io/)
- [Model Context Protocol](https://modelcontextprotocol.io/)

## 📄 License

This project is licensed under the ISC License - see the package.json file for details.

---

**Note**: This is a demonstration project. For production use, implement proper authentication, error handling, and security measures according to your requirements.