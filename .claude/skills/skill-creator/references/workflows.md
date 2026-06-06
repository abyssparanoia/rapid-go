# Workflow Patterns

## Sequential Workflow

Break complex tasks into clear sequential steps. Presenting a process overview at the top of SKILL.md is effective:

```markdown
Fill in the PDF form using the following steps:

1. Analyze the form (run analyze_form.py)
2. Create the field mapping (edit fields.json)
3. Validate the mapping (run validate_fields.py)
4. Fill in the form (run fill_form.py)
5. Verify the output (run verify_output.py)
```

## Conditional Branching Workflow

For tasks with branching logic, guide Claude through decision points:

```markdown
1. Determine the type of change:
   **Creating new content?** → Follow the "Creation Workflow"
   **Editing existing content?** → Follow the "Editing Workflow"

2. Creation Workflow: [steps]
3. Editing Workflow: [steps]
```
