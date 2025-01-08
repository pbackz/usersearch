import { NextResponse } from "next/server";
import { exec } from "child_process";
import { promisify } from "util";
import { parseGosearchResults } from "@/lib/parse-results";

const execAsync = promisify(exec);

// Mark this route as dynamic
export const dynamic = 'force-dynamic';

export async function POST(req: Request) {
  try {
    const { query } = await req.json();
    
    // Execute gosearch binary with the query
    const { stdout, stderr } = await execAsync(`bin/gosearch ${query}`);
    
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