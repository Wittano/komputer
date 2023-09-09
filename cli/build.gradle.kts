import com.github.jengelman.gradle.plugins.shadow.tasks.ShadowJar
import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

plugins {
    kotlin("jvm") version "1.9.10"
    kotlin("kapt") version "1.9.10"
    id("com.github.johnrengelman.shadow") version "8.1.1"

    application
}

group = "com.wittano.komputer.cli"
version = rootProject.version

repositories {
    mavenCentral()
}

val picocliVersion = findProperty("picocli.version") as? String

dependencies {
    // Internal dependencies
    implementation(project(":bot"))
    implementation(project(":commons"))

    // Kotlin
    implementation(kotlin("stdlib"))

    // Discord4j
    implementation("com.discord4j:discord4j-core:${findProperty("discord4j.version")}")
    implementation("io.projectreactor.kotlin:reactor-kotlin-extensions:1.2.2")

    // Picocli
    implementation("info.picocli:picocli:$picocliVersion")
    kapt("info.picocli:picocli-codegen:$picocliVersion")

    // Logger
    implementation("ch.qos.logback:logback-classic:${findProperty("logback-classic.version")}")
    implementation("org.codehaus.janino:janino:${findProperty("janino.version")}")

    testImplementation("org.junit.jupiter:junit-jupiter")
}

// TODO Create optional native image with GrallVM
// TODO Generate script to run CLI

tasks.test {
    useJUnitPlatform()
}

application {
    mainClass.set("com.wittano.komputer.cli.MainKt")
}

tasks.withType<ShadowJar> {
    manifest {
        attributes["Main-Class"] = application.mainClass
    }

    archiveBaseName.set("komputer-cli")
    archiveClassifier.set("")
}

tasks.withType<KotlinCompile> {
    kotlinOptions.jvmTarget = "17"
}