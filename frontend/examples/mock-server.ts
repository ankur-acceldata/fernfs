import { RemoteFS } from '../src';
import { MockWebSocketServer } from '../src/test/MockWebSocketServer';

// Start mock server
const server = new MockWebSocketServer();

// Create RemoteFS instance
const fs = new RemoteFS('ws://localhost:8080');

async function example() {
  try {
    // Create a directory
    await fs.mkdir('/documents');
    console.log('Created /documents directory');

    // Write a file
    await fs.writeFile('/documents/hello.txt', 'Hello, World!');
    console.log('Wrote to /documents/hello.txt');

    // Read the file
    const content = await fs.readFile('/documents/hello.txt', 'utf-8');
    console.log('File content:', content);

    // List directory contents
    const entries = await fs.readdir('/documents');
    console.log('Directory contents:', entries);

    // Get file stats
    const stats = await fs.stat('/documents/hello.txt');
    console.log('File stats:', stats);

    // Rename the file
    await fs.rename('/documents/hello.txt', '/documents/renamed.txt');
    console.log('Renamed file to /documents/renamed.txt');

    // Delete the file
    await fs.unlink('/documents/renamed.txt');
    console.log('Deleted /documents/renamed.txt');

    // Remove the directory
    await fs.rmdir('/documents');
    console.log('Removed /documents directory');

  } catch (error) {
    console.error('Error:', error);
  } finally {
    // Clean up
    await fs.destroy();
    server.close();
  }
}

// Run the example
example(); 