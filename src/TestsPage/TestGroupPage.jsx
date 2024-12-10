import { useState, useEffect } from "react";
import { useParams } from "react-router-dom";

export default function TestGroupPage() {
    const { id } = useParams();
    
    return (
        <section className="section">
            <h2 className="section__title"></h2>
        </section>
    )
}