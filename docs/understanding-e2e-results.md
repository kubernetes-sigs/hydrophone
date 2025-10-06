# Understanding E2E Test Results

This guide explains how to read and interpret the output from Kubernetes e2e conformance tests run through Hydrophone.

## Test Output Symbols

When tests are running, you'll see various symbols in the output that indicate test status:

### Symbol Meanings

- **`•`** (bullet point) = **Test Completed Successfully**
  - Each bullet represents one test that has finished running
  - Used for both passed and failed tests (check final results for pass/fail status)

- **`S`** = **Test Skipped**
  - The test was filtered out and didn't run
  - Common when using `--focus` or `--skip` flags to run specific subsets

- **`F`** = **Test Failed** 
  - The test ran but failed
  - Details will be shown in the final summary

- **`P`** = **Test Pending**
  - The test is marked as pending (not implemented or temporarily disabled)

### Example Output
```
SS•SS•SSS•SSSSSSSSSSSSSSS•SSS
```
This shows:
- 2 tests skipped (`SS`)
- 1 test completed (`•`)
- 2 tests skipped (`SS`)
- 1 test completed (`•`)
- 3 tests skipped (`SSS`)
- 1 test completed (`•`)
- Many more skipped tests...

## Understanding Test Counts

### "Will run X of Y specs"

When you see output like:
```
Will run 5 of 7392 specs
```

This means:
- **7392** = Total number of test specifications available in the e2e suite
- **5** = Number of tests that will actually run based on your filters

The large difference occurs because:
1. **Focus/Skip Filters**: Using `--focus` or `--skip` narrows down which tests run
2. **Conformance Tests**: Only ~400-500 tests are marked as `[Conformance]` out of 7000+ total
3. **Test Categories**: Many tests are marked as `[Serial]`, `[Slow]`, `[Disruptive]`, etc. and are filtered out

### Common Scenarios

| Scenario | Typical Count | Explanation |
|----------|---------------|-------------|
| Full conformance suite | ~400-500 of 7392 | Only `[Conformance]` tagged tests |
| Specific test focus | 1-50 of 7392 | Matching your `--focus` pattern |
| Custom filter | Varies | Based on your `--skip` and `--focus` patterns |

## Progress Indicators

Hydrophone shows real-time progress updates:

```
[15:04:05] Progress: 125/400 tests completed (31.3%)
```

- **Timestamp**: When the progress was checked
- **Completed/Total**: Number of tests finished vs. total to run
- **Percentage**: Progress completion

## Final Results Summary

At the end of a test run, you'll see a summary like:

```
Ran 400 of 7392 Specs in 45.123 seconds
SUCCESS! -- 398 Passed | 0 Failed | 2 Pending | 6992 Skipped
```

### Understanding the Summary
- **Ran**: Tests that actually executed
- **Passed**: Tests that completed successfully 
- **Failed**: Tests that encountered errors 
- **Pending**: Tests marked as not yet implemented
- **Skipped**: Tests filtered out by focus/skip patterns

## CNCF Conformance Submission

If you're running tests for CNCF Kubernetes conformance certification, you need:

### Required Files
1. **`e2e.log`** - Complete test output log
2. **`junit_01.xml`** - Machine-readable test results  
3. **`README.md`** - Instructions for reproducing results
4. **`PRODUCT.yaml`** - Product metadata

### Submission Process
1. **Run Conformance Tests**:
   ```bash
   hydrophone --conformance
   ```

2. **Collect Results**: Find output files in your specified `--output-dir`

3. **Submit to CNCF**: Create a pull request to [cncf/k8s-conformance](https://github.com/cncf/k8s-conformance)

### Submission Requirements
- All conformance tests must pass (no failures allowed)
- Use current or two most recent Kubernetes versions
- Include all required files in correct directory structure: `vX.Y/your-product-name/`

For detailed submission instructions, see: [CNCF K8s Conformance Instructions](https://github.com/cncf/k8s-conformance/blob/master/instructions.md)

## Troubleshooting Common Issues

### High Skip Count
This is normal! Most tests are filtered out for conformance runs. Only tests tagged with `[Conformance]` are required for certification.

### Long Runtime
Conformance tests typically take 1-2 hours to complete, depending on cluster performance and parallel execution settings.

## Configuration Options

Control test execution with these flags:

- `--focus 'pattern'` - Run only tests matching the pattern
- `--skip 'pattern'` - Skip tests matching the pattern  
- `--parallel N` - Run N tests in parallel (default: 1)
- `--verbosity N` - Set log verbosity level (0-4)
- `--disable-progress-status` - Turn off progress updates

## Getting Help

For questions about test results or conformance submission:
- [Kubernetes Slack #k8s-conformance](https://kubernetes.slack.com/channels/k8s-conformance)
- [Kubernetes Slack #sig-testing](https://kubernetes.slack.com/channels/sig-testing)  
- Email: conformance@cncf.io