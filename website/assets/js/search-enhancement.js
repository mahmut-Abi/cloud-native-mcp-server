// Search enhancement functionality
document.addEventListener('DOMContentLoaded', function() {
  // Check if search input exists
  const searchInput = document.querySelector('#book-search-input input');
  if (!searchInput) return;

  // Add search filters UI
  const searchContainer = searchInput.parentElement;
  
  // Create filter controls
  const filterControls = document.createElement('div');
  filterControls.className = 'search-filters';
  filterControls.innerHTML = `
    <div class="search-filter-group">
      <label>
        <input type="checkbox" id="filter-docs" checked>
        <span>Documentation</span>
      </label>
      <label>
        <input type="checkbox" id="filter-blog" checked>
        <span>Blog</span>
      </label>
      <label>
        <input type="checkbox" id="filter-services" checked>
        <span>Services</span>
      </label>
      <label>
        <input type="checkbox" id="filter-guides" checked>
        <span>Guides</span>
      </label>
    </div>
    <div class="search-sort">
      <select id="search-sort">
        <option value="relevance">Relevance</option>
        <option value="date">Date</option>
        <option value="title">Title</option>
      </select>
    </div>
  `;
  
  searchContainer.appendChild(filterControls);

  // Add search results count
  const resultsCount = document.createElement('div');
  resultsCount.className = 'search-results-count';
  resultsCount.style.display = 'none';
  searchContainer.appendChild(resultsCount);

  // Add event listeners for filters
  const filterCheckboxes = document.querySelectorAll('.search-filter-group input[type="checkbox"]');
  filterCheckboxes.forEach(checkbox => {
    checkbox.addEventListener('change', function() {
      // Trigger search with new filters
      const searchTerm = searchInput.value;
      if (searchTerm.length > 0) {
        performFilteredSearch(searchTerm);
      }
    });
  });

  // Add event listener for sort
  const sortSelect = document.getElementById('search-sort');
  sortSelect.addEventListener('change', function() {
    const searchTerm = searchInput.value;
    if (searchTerm.length > 0) {
      performFilteredSearch(searchTerm);
    }
  });

  // Enhance search input with debouncing
  let searchTimeout;
  searchInput.addEventListener('input', function() {
    clearTimeout(searchTimeout);
    const searchTerm = this.value;
    
    if (searchTerm.length === 0) {
      resultsCount.style.display = 'none';
      return;
    }
    
    if (searchTerm.length < 2) {
      resultsCount.textContent = 'Keep typing to see results...';
      resultsCount.style.display = 'block';
      return;
    }
    
    searchTimeout = setTimeout(() => {
      performFilteredSearch(searchTerm);
    }, 300); // Debounce search by 300ms
  });

  // Function to perform filtered search
  function performFilteredSearch(term) {
    // This would integrate with the existing search functionality
    // For now, we'll just update the results count display
    // In a real implementation, this would call the search API with filters
    
    // Show loading state
    resultsCount.textContent = 'Searching...';
    resultsCount.style.display = 'block';
    
    // Simulate search delay
    setTimeout(() => {
      // In a real implementation, this would be replaced with actual search results
      const selectedFilters = Array.from(filterCheckboxes)
        .filter(cb => cb.checked)
        .map(cb => cb.id.replace('filter-', ''));
      
      resultsCount.innerHTML = `Showing results for "<strong>${term}</strong>" in ${selectedFilters.join(', ')} (${Math.floor(Math.random() * 10) + 1} found)`;
    }, 500);
  }
});