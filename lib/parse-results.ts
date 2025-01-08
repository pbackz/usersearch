interface SearchResult {
    platform: string;
    url: string;
  }
  
  export function parseGosearchResults(output: string): SearchResult[] {
    const results: SearchResult[] = [];
    
    // Split output into lines and find lines with [32m[+]
    const lines = output.split('\n');
    const successPattern = /\[32m\[\+\] ([^:]+): (https?:\/\/[^\s]+)\[0m/;
    
    for (const line of lines) {
      const match = line.match(successPattern);
      if (match) {
        results.push({
          platform: match[1].trim(),
          url: match[2].trim()
        });
      }
    }
    
    return results;
  }