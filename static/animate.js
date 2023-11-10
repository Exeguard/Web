function animateElements(element, delay = 0) {
  let keyframes = [
    { transform: 'translateX(-200%)' },
    { transform: 'translateX(0)' }
  ];

  let options = {
    easing: 'ease-in-out',
    fill: 'forwards',
    duration: 1000,
    delay: delay,
  };

  element.animate(keyframes, options);

  Array.from(element.children).forEach((child, index) => {
    if (!child.classList.contains('card-body')) {
      animateElements(child, delay + 100 * (index + 1));
    }
  });
}

let mainElement = document.querySelector('main');
animateElements(mainElement);
