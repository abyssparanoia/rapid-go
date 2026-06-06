# Output Patterns

Patterns used in skills that require consistent, high-quality output.

## Template Pattern

Provide a template for the output format. Adjust the level of strictness to match the requirements.

**For strict requirements (API responses or data formats):**

```markdown
## Report Structure

Always use this exact template structure:

# [Analysis Title]

## Executive Summary
[One-paragraph overview of key findings]

## Key Findings
- Finding 1 backed by data
- Finding 2 backed by data
- Finding 3 backed by data

## Recommendations
1. Specific, actionable recommendation
2. Specific, actionable recommendation
```

**For flexible guidance (when adaptation is useful):**

```markdown
## Report Structure

The following is the default format, but adjust it using your best judgment:

# [Analysis Title]

## Executive Summary
[Overview]

## Key Findings
[Adapt sections based on what you find]

## Recommendations
[Tailor to the specific context]

Adjust sections based on the type of analysis.
```

## Example Pattern

For skills where output quality depends on seeing examples, provide input/output pairs:

```markdown
## Commit Message Format

Generate commit messages following these examples:

**Example 1:**
Input: Add user authentication with JWT tokens
Output:
feat(auth): implement JWT-based authentication

Add login endpoint and token validation middleware

**Example 2:**
Input: Fix bug where dates are not displayed correctly in reports
Output:
fix(reports): correct date formatting in timezone conversion

Use UTC timestamps consistently throughout report generation

Follow this style: type(scope): concise description, followed by a detailed explanation.
```

Examples communicate the desired style and level of detail to Claude more effectively than descriptions alone.
