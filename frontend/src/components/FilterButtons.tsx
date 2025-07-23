import React from 'react';

interface FilterButtonsProps {
  currentFilter: 'all' | 'active' | 'completed';
  onFilterChange: (filter: 'all' | 'active' | 'completed') => void;
}

const FilterButtons: React.FC<FilterButtonsProps> = ({ currentFilter, onFilterChange }) => {
  return (
    <div>
      <button
        onClick={() => onFilterChange('all')}
        style={{ fontWeight: currentFilter === 'all' ? 'bold' : 'normal' }}
      >
        All
      </button>
      <button
        onClick={() => onFilterChange('active')}
        style={{ fontWeight: currentFilter === 'active' ? 'bold' : 'normal' }}
      >
        Active
      </button>
      <button
        onClick={() => onFilterChange('completed')}
        style={{ fontWeight: currentFilter === 'completed' ? 'bold' : 'normal' }}
      >
        Completed
      </button>
    </div>
  );
};

export default FilterButtons;
