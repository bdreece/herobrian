export default class InviteLink extends HTMLInputElement {
    connectedCallback() {
        this.addEventListener('click', () => {
            this.select();
            document.execCommand('copy');
        });
    }
}

customElements.define('invite-link', InviteLink, { extends: 'input' });
