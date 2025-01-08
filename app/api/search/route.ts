import { NextResponse } from "next/server";
import { execFile } from "child_process";
import { promisify } from "util";
import { join } from "path";
import { parseGosearchResults } from "@/lib/parse-results";

const execFileAsync = promisify(execFile);

// Mark this route as dynamic
export const dynamic = 'force-dynamic';

export async function POST(req: Request) {
    try {
        const { query } = await req.json();
        
        // Use path.join for cross-platform compatibility
        const gosearchPath = join(process.cwd(), "bin", "gosearch");

        // Execute the gosearch binary with the provided query
        const { stdout, stderr } = await execFileAsync(gosearchPath, [query]);
        
        if (stderr) {
            return new NextResponse(stderr, { status: 500 });
        }

        // Parse the results
        const results = parseGosearchResults(stdout);
        return NextResponse.json(results);
    } catch (error) {
        console.error("Search failed:", error);
        return new NextResponse("Internal server error", { status: 500 });
    }
}
