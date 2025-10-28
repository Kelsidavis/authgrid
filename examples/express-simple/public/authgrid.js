/**
 * Authgrid Client SDK
 * Passwordless authentication using Ed25519 signatures
 */

class AuthgridClient {
  constructor(options = {}) {
    this.apiUrl = options.apiUrl || 'http://localhost:8080';
    this.storageKey = options.storageKey || 'authgrid_keypair';
  }

  /**
   * Register a new user and get a handle
   * @returns {Promise<{handle: string, publicKey: string}>}
   */
  async register() {
    try {
      // Generate Ed25519 keypair using WebCrypto API
      const keypair = await this.generateKeypair();

      // Export public key
      const publicKey = await this.exportPublicKey(keypair.publicKey);

      // Determine key type based on algorithm
      const keyType = keypair.privateKey.algorithm.name === 'Ed25519' ? 'ed25519' : 'ecdsa';

      // Send registration request
      const response = await fetch(`${this.apiUrl}/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          public_key: publicKey,
          key_type: keyType
        })
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Registration failed');
      }

      const data = await response.json();

      // Store keypair in localStorage (in production, use more secure storage)
      await this.storeKeypair(data.handle, keypair);

      return {
        handle: data.handle,
        publicKey: publicKey,
        id: data.id,
        createdAt: data.created_at
      };
    } catch (error) {
      console.error('Registration error:', error);
      throw error;
    }
  }

  /**
   * Authenticate with a handle
   * @param {string} handle - The user's handle
   * @returns {Promise<{token: string, expiresAt: string}>}
   */
  async authenticate(handle) {
    try {
      // Load keypair from storage
      const keypair = await this.loadKeypair(handle);
      if (!keypair) {
        throw new Error('Keypair not found. Please register first.');
      }

      // Request challenge
      const challengeResponse = await fetch(`${this.apiUrl}/challenge`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ handle })
      });

      if (!challengeResponse.ok) {
        const error = await challengeResponse.json();
        throw new Error(error.error || 'Challenge request failed');
      }

      const { challenge } = await challengeResponse.json();

      // Sign challenge
      const signature = await this.signChallenge(challenge, keypair.privateKey);

      // Verify signature
      const verifyResponse = await fetch(`${this.apiUrl}/verify`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          handle,
          challenge,
          signature
        })
      });

      if (!verifyResponse.ok) {
        const error = await verifyResponse.json();
        throw new Error(error.error || 'Verification failed');
      }

      const data = await verifyResponse.json();

      return {
        token: data.token,
        expiresAt: data.expires_at,
        verified: data.verified
      };
    } catch (error) {
      console.error('Authentication error:', error);
      throw error;
    }
  }

  /**
   * Generate Ed25519 keypair
   * @private
   */
  async generateKeypair() {
    // Note: Ed25519 is not directly supported in WebCrypto yet in all browsers
    // For MVP, we'll use a polyfill approach or ECDSA P-256 as fallback
    // In production, use a proper Ed25519 library like @noble/ed25519

    try {
      // Try to use Ed25519 (if available)
      return await crypto.subtle.generateKey(
        {
          name: 'Ed25519',
          namedCurve: 'Ed25519'
        },
        true,
        ['sign', 'verify']
      );
    } catch (e) {
      // Fallback to ECDSA P-256 for browsers without Ed25519
      console.warn('Ed25519 not supported, using ECDSA P-256 fallback');
      return await crypto.subtle.generateKey(
        {
          name: 'ECDSA',
          namedCurve: 'P-256'
        },
        true,
        ['sign', 'verify']
      );
    }
  }

  /**
   * Export public key to base64
   * @private
   */
  async exportPublicKey(publicKey) {
    const exported = await crypto.subtle.exportKey('spki', publicKey);
    return this.arrayBufferToBase64(exported);
  }

  /**
   * Sign challenge with private key
   * @private
   */
  async signChallenge(challengeBase64, privateKey) {
    const challengeBytes = this.base64ToArrayBuffer(challengeBase64);

    // Determine algorithm based on key type
    const keyAlgorithm = privateKey.algorithm.name;
    let signature;

    if (keyAlgorithm === 'Ed25519') {
      signature = await crypto.subtle.sign(
        'Ed25519',
        privateKey,
        challengeBytes
      );
    } else if (keyAlgorithm === 'ECDSA') {
      signature = await crypto.subtle.sign(
        {
          name: 'ECDSA',
          hash: { name: 'SHA-256' }
        },
        privateKey,
        challengeBytes
      );
    } else {
      throw new Error('Unsupported key algorithm');
    }

    return this.arrayBufferToBase64(signature);
  }

  /**
   * Store keypair in localStorage
   * @private
   */
  async storeKeypair(handle, keypair) {
    const publicKey = await crypto.subtle.exportKey('spki', keypair.publicKey);
    const privateKey = await crypto.subtle.exportKey('pkcs8', keypair.privateKey);

    const keyData = {
      handle,
      publicKey: this.arrayBufferToBase64(publicKey),
      privateKey: this.arrayBufferToBase64(privateKey),
      algorithm: keypair.privateKey.algorithm
    };

    localStorage.setItem(`${this.storageKey}_${handle}`, JSON.stringify(keyData));
  }

  /**
   * Load keypair from localStorage
   * @private
   */
  async loadKeypair(handle) {
    const stored = localStorage.getItem(`${this.storageKey}_${handle}`);
    if (!stored) return null;

    const keyData = JSON.parse(stored);

    const publicKey = await crypto.subtle.importKey(
      'spki',
      this.base64ToArrayBuffer(keyData.publicKey),
      keyData.algorithm,
      true,
      ['verify']
    );

    const privateKey = await crypto.subtle.importKey(
      'pkcs8',
      this.base64ToArrayBuffer(keyData.privateKey),
      keyData.algorithm,
      true,
      ['sign']
    );

    return { publicKey, privateKey };
  }

  /**
   * Convert ArrayBuffer to base64
   * @private
   */
  arrayBufferToBase64(buffer) {
    const bytes = new Uint8Array(buffer);
    let binary = '';
    for (let i = 0; i < bytes.byteLength; i++) {
      binary += String.fromCharCode(bytes[i]);
    }
    return btoa(binary);
  }

  /**
   * Convert base64 to ArrayBuffer
   * @private
   */
  base64ToArrayBuffer(base64) {
    const binary = atob(base64);
    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) {
      bytes[i] = binary.charCodeAt(i);
    }
    return bytes.buffer;
  }

  /**
   * Get stored handles
   * @returns {string[]}
   */
  getStoredHandles() {
    const handles = [];
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i);
      if (key && key.startsWith(this.storageKey + '_')) {
        const handle = key.replace(this.storageKey + '_', '');
        handles.push(handle);
      }
    }
    return handles;
  }

  /**
   * Remove stored keypair
   * @param {string} handle
   */
  removeKeypair(handle) {
    localStorage.removeItem(`${this.storageKey}_${handle}`);
  }
}

// Export for use in browser
if (typeof window !== 'undefined') {
  window.AuthgridClient = AuthgridClient;
}

// Export for Node.js/modules
if (typeof module !== 'undefined' && module.exports) {
  module.exports = AuthgridClient;
}
