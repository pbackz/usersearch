import { NextResponse } from "next/server";
import { execFile } from "child_process";
import { promisify } from "util";
import { join } from "path";
import { parseGosearchResults } from "@/lib/parse-results";
import { copyFileSync } from "fs";

const execFileAsync = promisify(execFile);

export const dynamic = 'force-dynamic';

export async function POST(req: Request) {
    try {
        const { query } = await req.json();
        
        // Prepare paths
        const gosearchPath = join(process.cwd(), "bin", "gosearch");
        const tmpGosearchPath = "/tmp/gosearch"; // Copy the binary to the writable tmp directory
        const tmpDataPath = "/tmp/data.json"; // Write data to the writable tmp directory
        
        // Optionally, copy a default data.json if needed
        copyFileSync(join(process.cwd(), "bin", "gosearch"), tmpGosearchPath);
        copyFileSync(join(process.cwd(), "bin", "data.json"), tmpDataPath);

        // Run the binary with the tmp directory as the target
        const { stdout, stderr } = await execFileAsync(tmpGosearchPath, [query, `--output=${tmpDataPath}`]);

        if (stderr) {
            return new NextResponse(stderr, { status: 500 });
        }
        
        // Read the results back from /tmp
        const results = parseGosearchResults(stdout);
        return NextResponse.json(results);
    } catch (error) {
        console.error("Search failed:", error);
        return new NextResponse("Internal server error", { status: 500 });
    }
}
