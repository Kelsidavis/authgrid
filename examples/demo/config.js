/**
 * Authgrid Frontend Configuration
 *
 * IMPORTANT: Update this file before deploying to production
 *
 * When hosting on https://authgrid.org/ with local backend:
 * 1. Expose your local API using ngrok, localtunnel, or Cloudflare Tunnel
 * 2. Update AUTHGRID_API_URL below with your public tunnel URL
 * 3. Deploy this frontend to your web host
 */

// ============================================================
// CONFIGURATION: Set your API endpoint URL here
// ============================================================

/**
 * Your backend API URL
 *
 * Options:
 * - Fly.io (deployed): 'https://authgrid-api.fly.dev'
 * - Local development: 'http://localhost:8080'
 * - Ngrok tunnel: 'https://abc123.ngrok-free.app'
 * - Localtunnel: 'https://your-subdomain.loca.lt'
 * - Cloudflare Tunnel: 'https://your-tunnel.trycloudflare.com'
 * - Your local IP (LAN only): 'http://192.168.1.100:8080'
 */
window.AUTHGRID_API_URL = 'https://authgrid-api.fly.dev';

// Example configurations (uncomment and update as needed):
// window.AUTHGRID_API_URL = 'https://abc123.ngrok-free.app';
// window.AUTHGRID_API_URL = 'https://authgrid-tunnel.trycloudflare.com';
