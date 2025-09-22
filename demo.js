// OpenTDF End-to-End Demo Implementation
// This demo shows the structure and flow of OpenTDF encryption/decryption
// For full functionality, connect to a real OpenTDF service

class OpenTDFDemo {
    constructor() {
        this.client = null;
        this.config = {
            kasEndpoint: 'https://demo.kas.opentdf.io',
            oidcOrigin: 'https://demo.keycloak.opentdf.io',
            clientId: 'demo-client',
            refreshToken: ''
        };
        this.loadConfig();
        this.initializeSDK();
    }

    async initializeSDK() {
        // Check if we're in a module-supported browser
        try {
            // Try to dynamically import the OpenTDF SDK
            const { OpenTDF, AuthProviders } = await import('@opentdf/sdk');
            this.OpenTDF = OpenTDF;
            this.AuthProviders = AuthProviders;
            this.sdkAvailable = true;
            console.log('OpenTDF SDK loaded successfully');
        } catch (error) {
            console.warn('OpenTDF SDK not available in this environment, using demo mode:', error.message);
            this.sdkAvailable = false;
            // In demo mode, we'll simulate the functionality
        }
    }

    loadConfig() {
        const savedConfig = localStorage.getItem('opentdf-config');
        if (savedConfig) {
            this.config = { ...this.config, ...JSON.parse(savedConfig) };
            this.updateConfigInputs();
        }
    }

    saveConfig() {
        this.config = {
            kasEndpoint: document.getElementById('kasEndpoint').value,
            oidcOrigin: document.getElementById('oidcOrigin').value,
            clientId: document.getElementById('clientId').value,
            refreshToken: document.getElementById('refreshToken').value
        };
        localStorage.setItem('opentdf-config', JSON.stringify(this.config));
        this.showStatus('Configuration saved successfully!', 'success', 'encryptStatus');
    }

    updateConfigInputs() {
        document.getElementById('kasEndpoint').value = this.config.kasEndpoint;
        document.getElementById('oidcOrigin').value = this.config.oidcOrigin;
        document.getElementById('clientId').value = this.config.clientId;
        document.getElementById('refreshToken').value = this.config.refreshToken;
    }

    async initializeClient() {
        try {
            if (!this.sdkAvailable) {
                this.showStatus('Demo mode: Using simulated OpenTDF client', 'info', 'encryptStatus');
                return true;
            }

            // For demo purposes, we'll create a simple auth provider
            // In a real application, you would use proper OIDC authentication
            let authProvider;
            
            if (this.config.refreshToken) {
                const oidcCredentials = {
                    clientId: this.config.clientId,
                    exchange: 'refresh',
                    refreshToken: this.config.refreshToken,
                    oidcOrigin: this.config.oidcOrigin,
                };
                authProvider = await this.AuthProviders.refreshAuthProvider(oidcCredentials);
            } else {
                // For demo without real auth, create a basic auth provider
                authProvider = await this.AuthProviders.clientSecretAuthProvider({
                    clientId: this.config.clientId,
                    clientSecret: 'demo-secret', // This would be configured properly in real use
                    oidcOrigin: this.config.oidcOrigin,
                });
            }

            this.client = new this.OpenTDF({
                authProvider,
                defaultCreateOptions: {
                    defaultKASEndpoint: this.config.kasEndpoint,
                },
                dpopKeys: authProvider.getSigningKey(),
            });

            return true;
        } catch (error) {
            console.error('Failed to initialize OpenTDF client:', error);
            this.showStatus(`Failed to initialize client: ${error.message}`, 'error', 'encryptStatus');
            return false;
        }
    }

    async encryptFile() {
        const fileInput = document.getElementById('fileToEncrypt');
        const file = fileInput.files[0];
        
        if (!file) {
            this.showStatus('Please select a file to encrypt', 'error', 'encryptStatus');
            return;
        }

        try {
            this.showStatus('Initializing OpenTDF client...', 'info', 'encryptStatus');
            const clientInitialized = await this.initializeClient();
            
            if (!clientInitialized) {
                return;
            }

            this.showStatus('Encrypting file...', 'info', 'encryptStatus');
            
            if (this.sdkAvailable && this.client) {
                // Real encryption using OpenTDF SDK
                const source = new ReadableStream({
                    start(controller) {
                        const reader = new FileReader();
                        reader.onload = () => {
                            const uint8Array = new Uint8Array(reader.result);
                            controller.enqueue(uint8Array);
                            controller.close();
                        };
                        reader.onerror = () => controller.error(reader.error);
                        reader.readAsArrayBuffer(file);
                    }
                });

                // Create encrypted TDF
                const encryptedStream = await this.client.createNanoTDF({
                    source: { type: 'stream', location: source },
                });

                // Convert stream back to blob for download
                const encryptedData = await this.streamToBlob(encryptedStream);
                
                this.showStatus('File encrypted successfully!', 'success', 'encryptStatus');
                this.displayEncryptedResult(encryptedData, file.name);
            } else {
                // Demo mode: simulate encryption
                await this.simulateEncryption(file);
            }

        } catch (error) {
            console.error('Encryption failed:', error);
            this.showStatus(`Encryption failed: ${error.message}`, 'error', 'encryptStatus');
        }
    }

    async simulateEncryption(file) {
        // Simulate encryption process for demo
        await this.delay(2000); // Simulate processing time
        
        const reader = new FileReader();
        reader.onload = () => {
            const originalData = new Uint8Array(reader.result);
            
            // Create a simulated encrypted file with TDF header
            const tdfHeader = new TextEncoder().encode('TDF_DEMO_ENCRYPTED_');
            const metadata = new TextEncoder().encode(JSON.stringify({
                filename: file.name,
                size: file.size,
                type: file.type,
                encrypted: true,
                timestamp: Date.now(),
                format: 'NanoTDF'
            }));
            
            // Simple simulation: reverse the bytes (not real encryption!)
            const reversedData = new Uint8Array(originalData.length);
            for (let i = 0; i < originalData.length; i++) {
                reversedData[i] = originalData[originalData.length - 1 - i];
            }
            
            // Combine header + metadata + "encrypted" data
            const totalSize = tdfHeader.length + 4 + metadata.length + reversedData.length;
            const encryptedData = new Uint8Array(totalSize);
            let offset = 0;
            
            encryptedData.set(tdfHeader, offset);
            offset += tdfHeader.length;
            
            // Add metadata length (4 bytes)
            const metadataLength = new Uint32Array([metadata.length]);
            encryptedData.set(new Uint8Array(metadataLength.buffer), offset);
            offset += 4;
            
            encryptedData.set(metadata, offset);
            offset += metadata.length;
            
            encryptedData.set(reversedData, offset);
            
            const encryptedBlob = new Blob([encryptedData], { type: 'application/octet-stream' });
            
            this.showStatus('File encrypted successfully! (Demo Mode)', 'success', 'encryptStatus');
            this.displayEncryptedResult(encryptedBlob, file.name);
        };
        reader.readAsArrayBuffer(file);
    }

    async decryptFile() {
        const fileInput = document.getElementById('fileToDecrypt');
        const file = fileInput.files[0];
        
        if (!file) {
            this.showStatus('Please select an encrypted file to decrypt', 'error', 'decryptStatus');
            return;
        }

        try {
            this.showStatus('Initializing OpenTDF client...', 'info', 'decryptStatus');
            const clientInitialized = await this.initializeClient();
            
            if (!clientInitialized) {
                return;
            }

            this.showStatus('Decrypting file...', 'info', 'decryptStatus');
            
            if (this.sdkAvailable && this.client) {
                // Real decryption using OpenTDF SDK
                const encryptedStream = new ReadableStream({
                    start(controller) {
                        const reader = new FileReader();
                        reader.onload = () => {
                            const uint8Array = new Uint8Array(reader.result);
                            controller.enqueue(uint8Array);
                            controller.close();
                        };
                        reader.onerror = () => controller.error(reader.error);
                        reader.readAsArrayBuffer(file);
                    }
                });

                // Decrypt the TDF
                const decryptedStream = await this.client.read({
                    type: 'stream',
                    location: encryptedStream
                });

                // Convert stream back to text/blob for display
                const decryptedData = await this.streamToText(decryptedStream);
                
                this.showStatus('File decrypted successfully!', 'success', 'decryptStatus');
                this.displayDecryptedResult(decryptedData, file.name);
            } else {
                // Demo mode: simulate decryption
                await this.simulateDecryption(file);
            }

        } catch (error) {
            console.error('Decryption failed:', error);
            this.showStatus(`Decryption failed: ${error.message}`, 'error', 'decryptStatus');
        }
    }

    async simulateDecryption(file) {
        // Simulate decryption process for demo
        await this.delay(1500); // Simulate processing time
        
        const reader = new FileReader();
        reader.onload = () => {
            const encryptedData = new Uint8Array(reader.result);
            
            try {
                // Check for demo TDF header
                const headerText = new TextDecoder().decode(encryptedData.slice(0, 18));
                if (!headerText.startsWith('TDF_DEMO_ENCRYPTED_')) {
                    throw new Error('Not a valid demo TDF file');
                }
                
                // Extract metadata
                let offset = 18;
                const metadataLength = new Uint32Array(encryptedData.slice(offset, offset + 4).buffer)[0];
                offset += 4;
                
                const metadataBytes = encryptedData.slice(offset, offset + metadataLength);
                const metadata = JSON.parse(new TextDecoder().decode(metadataBytes));
                offset += metadataLength;
                
                // Extract and "decrypt" the data (reverse the bytes back)
                const reversedData = encryptedData.slice(offset);
                const originalData = new Uint8Array(reversedData.length);
                for (let i = 0; i < reversedData.length; i++) {
                    originalData[i] = reversedData[reversedData.length - 1 - i];
                }
                
                // Convert to text for display
                let decryptedContent;
                try {
                    decryptedContent = new TextDecoder().decode(originalData);
                } catch (e) {
                    // If it's binary data, show hex representation
                    decryptedContent = Array.from(originalData)
                        .map(b => b.toString(16).padStart(2, '0'))
                        .join(' ');
                    decryptedContent = `Binary data (${originalData.length} bytes):\n${decryptedContent}`;
                }
                
                this.showStatus('File decrypted successfully! (Demo Mode)', 'success', 'decryptStatus');
                this.displayDecryptedResult(decryptedContent, metadata.filename || file.name);
                
            } catch (error) {
                throw new Error('Invalid or corrupted TDF file');
            }
        };
        reader.readAsArrayBuffer(file);
    }

    delay(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }

    async streamToBlob(stream) {
        const reader = stream.getReader();
        const chunks = [];
        
        while (true) {
            const { done, value } = await reader.read();
            if (done) break;
            chunks.push(value);
        }
        
        return new Blob(chunks);
    }

    async streamToText(stream) {
        const reader = stream.getReader();
        const decoder = new TextDecoder();
        let result = '';
        
        while (true) {
            const { done, value } = await reader.read();
            if (done) break;
            result += decoder.decode(value, { stream: true });
        }
        
        return result;
    }

    displayEncryptedResult(encryptedBlob, originalFileName) {
        const resultDiv = document.getElementById('encryptResult');
        const downloadUrl = URL.createObjectURL(encryptedBlob);
        const encryptedFileName = `${originalFileName}.ntdf`;
        
        resultDiv.innerHTML = `
            <h4>✅ Encryption Complete</h4>
            <p><strong>Original file:</strong> ${originalFileName}</p>
            <p><strong>Encrypted file size:</strong> ${this.formatFileSize(encryptedBlob.size)}</p>
            <p><strong>Format:</strong> NanoTDF ${this.sdkAvailable ? '' : '(Demo Mode)'}</p>
            <a href="${downloadUrl}" download="${encryptedFileName}" 
               style="display: inline-block; background-color: #27ae60; color: white; 
                      padding: 10px 20px; text-decoration: none; border-radius: 5px; margin-top: 10px;">
                📥 Download Encrypted File
            </a>
        `;
    }

    displayDecryptedResult(decryptedText, encryptedFileName) {
        const resultDiv = document.getElementById('decryptResult');
        const originalFileName = encryptedFileName.replace('.ntdf', '').replace('.tdf', '');
        
        resultDiv.innerHTML = `
            <h4>✅ Decryption Complete</h4>
            <p><strong>Encrypted file:</strong> ${encryptedFileName}</p>
            <p><strong>Decrypted content size:</strong> ${this.formatFileSize(new Blob([decryptedText]).size)}</p>
            <p><strong>Mode:</strong> ${this.sdkAvailable ? 'Real OpenTDF' : 'Demo Simulation'}</p>
            <div style="background-color: #f8f9fa; padding: 15px; border-radius: 5px; 
                        border: 1px solid #e9ecef; margin-top: 10px; max-height: 300px; overflow-y: auto;">
                <h5>📄 Decrypted Content:</h5>
                <pre style="white-space: pre-wrap; word-wrap: break-word; margin: 0;">${this.escapeHtml(decryptedText)}</pre>
            </div>
            <button onclick="demo.downloadDecryptedContent('${this.escapeHtml(decryptedText)}', '${originalFileName}')" 
                    style="background-color: #27ae60; color: white; padding: 10px 20px; 
                           border: none; border-radius: 5px; margin-top: 10px; cursor: pointer;">
                📥 Download Decrypted File
            </button>
        `;
    }

    downloadDecryptedContent(content, filename) {
        const blob = new Blob([content], { type: 'text/plain' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    }

    createSampleTextFile() {
        const content = `Hello from OpenTDF Demo!

This is a sample text file created for testing OpenTDF encryption and decryption.

Created on: ${new Date().toISOString()}

You can encrypt this file using the demo above, and then decrypt it to see the original content.

Features demonstrated:
- File encryption using OpenTDF NanoTDF format
- Secure decryption with proper authentication
- End-to-end encryption workflow
- Web-based interface for easy testing

OpenTDF provides zero-trust data protection that ensures your sensitive data remains secure throughout its lifecycle.

Mode: ${this.sdkAvailable ? 'Real OpenTDF SDK Integration' : 'Demo Simulation Mode'}`;

        this.downloadSampleFile(content, 'sample-text.txt', 'text/plain');
    }

    createSampleJsonFile() {
        const content = {
            "demo": "OpenTDF End-to-End Encryption",
            "timestamp": new Date().toISOString(),
            "sdk_available": this.sdkAvailable,
            "data": {
                "message": "This is sensitive data that should be encrypted",
                "users": ["alice", "bob", "charlie"],
                "permissions": {
                    "read": true,
                    "write": false,
                    "share": true
                }
            },
            "metadata": {
                "classification": "confidential",
                "retention_days": 365,
                "geographic_restriction": ["US", "EU"]
            }
        };

        this.downloadSampleFile(JSON.stringify(content, null, 2), 'sample-data.json', 'application/json');
    }

    downloadSampleFile(content, filename, mimeType) {
        const blob = new Blob([content], { type: mimeType });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
        
        this.showStatus(`Sample file "${filename}" downloaded successfully!`, 'success', 'encryptStatus');
    }

    showStatus(message, type, elementId) {
        const statusElement = document.getElementById(elementId);
        statusElement.className = `status ${type}`;
        statusElement.textContent = message;
        statusElement.style.display = 'block';
        
        if (type === 'success') {
            setTimeout(() => {
                statusElement.style.display = 'none';
            }, 5000);
        }
    }

    formatFileSize(bytes) {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

// Initialize the demo
const demo = new OpenTDFDemo();

// Make functions available globally for HTML onclick handlers
window.saveConfig = () => demo.saveConfig();
window.encryptFile = () => demo.encryptFile();
window.decryptFile = () => demo.decryptFile();
window.createSampleTextFile = () => demo.createSampleTextFile();
window.createSampleJsonFile = () => demo.createSampleJsonFile();
window.demo = demo; // For accessing downloadDecryptedContent