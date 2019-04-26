class App {
    constructor() {
        this.account = Orbs.createAccount();
        this.client = new Orbs.Client("http://localhost:8080", 42, Orbs.NetworkType.NETWORK_TYPE_MAIN_NET);
    }

    getNameFromURL() {
        return decodeURIComponent(window.location.search.substr(1) || "Main");
    }

    async getLastRevision() {
        const query = this.client.createQuery(this.account.publicKey, "Lapti", "GetLastRevision",
            [Orbs.argString(this.getNameFromURL())]);
        const response = await this.client.sendQuery(query);
        this.revision = JSON.parse(response.outputArguments[0].value);
        if (!this.revision.Name) {
            this.revision.Name = this.getNameFromURL();
        }

        return this.revision;
    }

    async saveRevision() {
        const [ tx ] = this.client.createTransaction(this.account.publicKey, this.account.privateKey,
            "Lapti", "SaveRevision", [Orbs.argString(this.revision.Name), Orbs.argString(this.revision.Text)]);
        const response = await this.client.sendTransaction(tx);
        console.log(response);
        this.revision = JSON.parse(response.outputArguments[0].value);

        return this.revision;
    }

    async onLoad() {
        await this.getLastRevision();
        await this.render();
    }

    async render() {
        const converter = new showdown.Converter();

        this.updateElement("article_text", converter.makeHtml(this.revision.Text));
        this.updateElement("article_name", this.revision.Name);
        this.updateElement("article_content", this.revision.Text);

        this.updateElement("public_key", this.account.address);
    }

    updateElement(id, value) {
        const element = document.getElementById(id);
        if (id) {
            element.innerHTML = value;
        }
    }

    async submitForm() {
        const textValue = document.getElementById('article_content').value;
        this.revision.Text = textValue;

        await this.saveRevision();
        await this.render();
        this.hideEditForm();

        return false;
    }

    showEditForm() {
        document.getElementById("article_text").classList.add("invisible");
        document.getElementById("edit_form").classList.remove("invisible");
    }

    hideEditForm() {
        document.getElementById("article_text").classList.remove("invisible");
        document.getElementById("edit_form").classList.add("invisible");
    }
}