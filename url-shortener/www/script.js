// Wait for the DOM to be fully loaded
document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM loaded, initializing URL shortener...');
    
    // Get DOM elements
    const form = document.getElementById('urlForm');
    const urlInput = document.getElementById('urlInput');
    const shortenBtn = document.getElementById('shortenBtn');
    const loading = document.getElementById('loading');
    const result = document.getElementById('result');
    const shortenedUrl = document.getElementById('shortenedUrl');
    const copyBtn = document.getElementById('copyBtn');
    const errorMessage = document.getElementById('errorMessage');
    
    // Debug: Check if elements are found
    console.log('Form found:', !!form);
    console.log('URL input found:', !!urlInput);
    console.log('Button found:', !!shortenBtn);
    
    if (!form || !urlInput || !shortenBtn) {
        console.error('Required elements not found!');
        return;
    }
    
    // Helper functions
    function isValidUrl(string) {
        try {
            new URL(string);
            return true;
        } catch (_) {
            return false;
        }
    }
    
    function setLoading(isLoading) {
        if (isLoading) {
            loading.style.display = 'block';
            shortenBtn.disabled = true;
            shortenBtn.textContent = 'Shortening...';
        } else {
            loading.style.display = 'none';
            shortenBtn.disabled = false;
            shortenBtn.textContent = 'Shorten URL';
        }
    }
    
    function showResult(type) {
        result.className = `result ${type}`;
        result.style.display = 'block';
    }
    
    function hideResult() {
        result.style.display = 'none';
    }
    
    function showError(message) {
        errorMessage.textContent = message;
        errorMessage.style.display = 'block';
        console.log('Error shown:', message);
    }
    
    function hideError() {
        errorMessage.style.display = 'none';
    }
    
    function copyToClipboard(text) {
        if (navigator.clipboard && window.isSecureContext) {
            // Modern approach
            navigator.clipboard.writeText(text).then(() => {
                showCopySuccess();
            }).catch(err => {
                console.error('Clipboard API failed:', err);
                fallbackCopyToClipboard(text);
            });
        } else {
            // Fallback for older browsers
            fallbackCopyToClipboard(text);
        }
    }
    
    function fallbackCopyToClipboard(text) {
        const textArea = document.createElement('textarea');
        textArea.value = text;
        textArea.style.position = 'fixed';
        textArea.style.left = '-999999px';
        textArea.style.top = '-999999px';
        document.body.appendChild(textArea);
        textArea.focus();
        textArea.select();
        
        try {
            document.execCommand('copy');
            showCopySuccess();
        } catch (err) {
            console.error('Fallback copy failed:', err);
        } finally {
            document.body.removeChild(textArea);
        }
    }
    
    function showCopySuccess() {
        const originalText = copyBtn.textContent;
        copyBtn.textContent = 'Copied!';
        copyBtn.style.background = '#059669';
        
        setTimeout(() => {
            copyBtn.textContent = originalText;
            copyBtn.style.background = '#10b981';
        }, 2000);
    }
    
    // Main submit function
    async function handleSubmit(e) {
        e.preventDefault();
        console.log('Form submitted');
        
        const url = urlInput.value.trim();
        console.log('URL entered:', url);
        
        // Validate input
        if (!url) {
            showError('Please enter a URL');
            return;
        }
        
        // Validate URL format
        if (!isValidUrl(url)) {
            showError('Please enter a valid URL (including http:// or https://)');
            return;
        }
        
        // Clear previous results
        hideError();
        hideResult();
        
        // Show loading
        setLoading(true);
        
        try {
            console.log('Sending request to server...');
            const response = await fetch('/', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ url: url })
            });
            
            console.log('Response status:', response.status);
            
            if (!response.ok) {
                throw new Error(`Server responded with status: ${response.status}`);
            }
            
            const data = await response.json();
            console.log('Response data:', data);
            
            // Display result
            const shortUrl = `${window.location.origin}/${data.ID}`;
            shortenedUrl.innerHTML = `<a href="${shortUrl}" target="_blank" rel="noopener noreferrer">${shortUrl}</a>`;
            showResult('success');
            
            // Setup copy button
            copyBtn.onclick = () => copyToClipboard(shortUrl);
            
        } catch (error) {
            console.error('Error shortening URL:', error);
            showError('Failed to shorten URL. Please try again.');
        } finally {
            setLoading(false);
        }
    }
    
    // Event listeners
    form.addEventListener('submit', handleSubmit);
    
    urlInput.addEventListener('input', function() {
        hideError();
        hideResult();
    });
    
    console.log('URL shortener initialized successfully');
});
