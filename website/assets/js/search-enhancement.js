// Search enhancement functionality
document.addEventListener('DOMContentLoaded', function() {
  
  // HTML escaping function to prevent XSS
  function escapeHTML(text) {
    if (!text) return '';
    
    // Use DOM API for safe HTML escaping
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  }
  // Check if search input exists
  const searchInput = document.querySelector('#book-search-input input');
  if (!searchInput) return;

  // Add search filters UI
  const searchContainer = searchInput.parentElement;
  
  // Create filter controls
  const filterControls = document.createElement('div');
  filterControls.className = 'search-filters';
  
  // Create filter group
  const filterGroup = document.createElement('div');
  filterGroup.className = 'search-filter-group';
  
  // Define filters configuration
  const filters = [
    {id: 'filter-docs', label: 'Documentation', checked: true},
    {id: 'filter-blog', label: 'Blog', checked: true},
    {id: 'filter-services', label: 'Services', checked: true},
    {id: 'filter-guides', label: 'Guides', checked: true}
  ];
  
  // Create filter elements programmatically
  filters.forEach(filter => {
    const label = document.createElement('label');
    const input = document.createElement('input');
    input.type = 'checkbox';
    input.id = filter.id;
    input.checked = filter.checked;
    
    const span = document.createElement('span');
    span.textContent = filter.label;
    
    label.appendChild(input);
    label.appendChild(span);
    filterGroup.appendChild(label);
  });
  
  // Create sort dropdown
  const sortDropdown = document.createElement('div');
  sortDropdown.className = 'search-sort';
  
  const select = document.createElement('select');
  select.id = 'search-sort';
  
  const options = [
    {value: 'relevance', text: 'Relevance'},
    {value: 'date', text: 'Date'},
    {value: 'title', text: 'Title'}
  ];
  
  options.forEach(option => {
    const optionElement = document.createElement('option');
    optionElement.value = option.value;
    optionElement.textContent = option.text;
    select.appendChild(optionElement);
  });
  
  sortDropdown.appendChild(select);
  
  // Assemble the controls
  filterControls.appendChild(filterGroup);
  filterControls.appendChild(sortDropdown);
  
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
      
      resultsCount.innerHTML = `Showing results for "<strong>${escapeHTML(term)}</strong>" in ${escapeHTML(selectedFilters.join(', '))} (${Math.floor(Math.random() * 10) + 1} found)`;
    }, 500);
  }
});