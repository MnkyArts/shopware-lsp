# Shopware LSP - Performance and Usage Guidelines

## Performance Optimizations

### Indexing Strategy
- **Incremental Updates**: Only changed files are re-indexed
- **Selective Indexing**: Files are filtered by extension and location before parsing
- **Batch Operations**: Multiple items are saved in single database transactions
- **Lazy Loading**: Component/entity details are loaded on-demand

### Memory Management
- **Tree Close**: Tree-sitter trees are properly closed after parsing to prevent memory leaks
- **Parser Pooling**: Parsers are created and destroyed per operation to avoid state issues
- **Bounded Results**: Completion results are limited to prevent overwhelming the client

### Database Optimization
- **BoltDB**: Fast key-value store with O(log n) lookup time
- **Indexed Keys**: Primary keys for fast lookups by name
- **Batch Writes**: Multiple items written in single transaction
- **File-based Deletion**: Efficient removal of all items from deleted files

## Usage Best Practices

### For Large Projects
1. **Initial Indexing**: First startup may take 30-60 seconds for large projects
2. **Incremental Updates**: Subsequent changes are indexed in milliseconds
3. **Memory Usage**: Expect 50-200MB RAM usage depending on project size

### Completion Performance
- **Context-Aware**: Providers only activate in relevant contexts
- **Filtered Results**: Results are filtered by relevance before returning
- **Trigger Characters**: Smart triggering reduces unnecessary completions

### Troubleshooting
If performance issues occur:
1. Check LSP logs for errors
2. Clear cache: Delete `.cache/shopware-lsp/` directory
3. Force reindex: Use "Shopware: Force Reindex" command
4. Restart LSP server: Use "Restart Shopware Language Server" command

## Known Limitations

### Phase 1-5 Implementation
- Abstract class diagnostics deferred (requires extensive PHP AST analysis)
- Code lens for block versioning deferred (requires Shopware core integration)
- Additional generators (scheduled task, changelog) can be added in future

### Performance Considerations
- Very large JavaScript files (>5000 lines) may have slower completion
- Projects with >10,000 components may see increased memory usage
- Initial indexing scales linearly with project size

## Future Optimizations (Phase 6+)
- Parallel indexing for multi-core systems
- Incremental parsing for large files
- Caching of frequently accessed completions
- Streaming results for large result sets
